package pokeAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/"

var (
	location           *Locations
	locationAreaList   *LocationAreaList
	locationAreaDetail *LocationAreaDetail
	cache              *pokecache.Cache
)

type Locations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaList struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
	GameIndices []struct {
		GameIndex  int `json:"game_index"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"game_indices"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}

type LocationAreaDetail struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func fetchFromAPI(url string, v interface{}) error {
	if cache == nil {
		cache = pokecache.NewCache(time.Minute * 5)
	}

	// Check if the URL's data is in the cache
	if data, found := cache.Get(url); found {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return nil
	}

	// Fetch data from API
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get request failed: %w", err)
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle HTTP status codes
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("bad HTTP response: status code: %d", res.StatusCode)
	}

	// Unmarshal the JSON response
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	// Add to cache
	cache.Add(url, body)
	return nil
}

func init() {
	cache = pokecache.NewCache(time.Minute * 5)
}

func initLocations() *Locations {
	location = new(Locations)
	location.fetch(baseURL + "location-area?offset=0&limit=20")
	return location
}

func GetLocations() *Locations {
	if location == nil {
		return initLocations()
	}
	return location
}

func (l *Locations) fetch(url string) {
	fetchFromAPI(url, l)
}

func (l *Locations) GetNext() error {
	if l.Next == nil {
		return errors.New("no further locations to get")
	}
	l.fetch(*l.Next)
	return nil
}

func (l *Locations) GetPrevious() error {
	if l.Previous == nil {
		return errors.New("at the start of locations list, cannot get previous page")
	}
	l.fetch(*l.Previous)
	return nil
}

func (l *Locations) PrintLocations() error {
	if len(l.Results) == 0 {
		return errors.New("no locations to print")
	}
	for _, location := range l.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func initLocationArea() *LocationAreaList {
	locationAreaList = new(LocationAreaList)
	return locationAreaList
}

func GetLocationArea() *LocationAreaList {
	if locationAreaList == nil {
		return initLocationArea()
	}
	return locationAreaList
}

func initLocationAreaDetail() *LocationAreaDetail {
	locationAreaDetail = new(LocationAreaDetail)
	return locationAreaDetail
}

func GetLocationAreaDetail() *LocationAreaDetail {
	if locationAreaDetail == nil {
		return initLocationAreaDetail()
	}
	return locationAreaDetail
}

func (l *LocationAreaDetail) GetLocationDetail(name string) error {
	url := baseURL + "location-area/" + name
	return fetchFromAPI(url, l)
}

func (l *LocationAreaDetail) GetLocationPokemon(name string) (pokemon []string, err error) {
	err = l.GetLocationDetail(name)
	if err != nil {
		return
	}
	for _, p := range l.PokemonEncounters {
		pokemon = append(pokemon, p.Pokemon.Name)
	}

	return
}
