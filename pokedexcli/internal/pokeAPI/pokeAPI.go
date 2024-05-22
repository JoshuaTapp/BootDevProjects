package pokeAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokecache"
)

var baseURL = "https://pokeapi.co/api/v2/"
var location *Locations
var cache *pokecache.Cache

type Locations struct {
	Count     int     `json:"count"`
	Next      *string `json:"next"`
	Previous  *string `json:"previous"`
	Locations []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Location struct {
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
	ID    int    `json:"id"`
	Name  string `json:"name"`
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

func InitLocations() *Locations {
	cache = pokecache.NewCache(time.Minute * 5)

	url := baseURL + "location?offset=0&limit=20"
	location = new(Locations)
	location.getNewData(url)
	return location
}

func GetLocations() *Locations {
	if location == nil {
		InitLocations()
	}
	return location
}

// Handles the HTTP request to get the locations and unmarshal the JSON response.
// Updates the receiver 'locations' struct with the response's data.
func (l *Locations) getNewData(url string) {
	if data, ok := cache.Get(url); ok {
		err := json.Unmarshal(data, l)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	//log.Default().Printf("Response body: %s\n", body)

	if res.StatusCode < http.StatusOK && res.StatusCode >= http.StatusMultipleChoices {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, l)
	if err != nil {
		log.Fatalf("Error: %v\nFailed to unmarshal with json payload: %v\n", err, l)
	} else {
		cache.Add(url, body)
	}
}

func (l *Locations) GetNext() error {
	url := l.Next
	if url == nil {
		log.Default().Println("l.next is nil, unable to get next page")
		return errors.New("no further new locations to get")
	}
	l.getNewData(*url)
	return nil
}

func (l *Locations) GetPrevious() error {
	url := l.Previous
	if url == nil {
		log.Default().Println("l.previous is nil, unable to get previous page")
		return errors.New("at start of locations list, cannot get previous page")
	}
	l.getNewData(*url)
	return nil
}

func (l *Locations) PrintLocations() error {
	if len(l.Locations) < 1 {
		log.Default().Println("l.locations is empty, unable to print locations")
		return errors.New("no locations to print")
	}
	for _, location := range l.Locations {
		fmt.Println(location.Name)
	}
	return nil
}
