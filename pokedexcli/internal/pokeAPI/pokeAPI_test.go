package pokeAPI

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokecache"
// 	. "github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
// )

// func setup() {
// 	cache = pokecache.NewCache(time.Minute * 5)
// }

// func teardown() {
// 	cache = nil
// 	location = nil
// 	locationAreaList = nil
// }

// func TestGetLocations(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	t.Run("Initial call to GetLocations", func(t *testing.T) {
// 		// Mock server and response for /location-area?offset=0&limit=20
// 		handler := func(w http.ResponseWriter, r *http.Request) {
// 			if strings.Contains(r.URL.Path, "location-area") {
// 				w.Header().Set("Content-Type", "application/json")
// 				w.Write([]byte(`{
// 					"count": 2,
// 					"next": null,
// 					"previous": null,
// 					"results": [
// 						{"name": "location1", "url": "url1"},
// 						{"name": "location2", "url": "url2"}
// 					]
// 				}`))
// 			}
// 		}
// 		server := httptest.NewServer(http.HandlerFunc(handler))
// 		defer server.Close()

// 		// Override baseURL for testing
// 		baseURL = server.URL + "/"

// 		loc := GetLocations()
// 		if loc == nil || len(loc.Locations) != 2 {
// 			t.Fatalf("expected 2 locations, got %v", loc)
// 		}
// 	})

// 	t.Run("GetLocations cached", func(t *testing.T) {
// 		cache.Add(baseURL+"location-area?offset=0&limit=20", []byte(`{
// 			"count": 2,
// 			"next": null,
// 			"previous": null,
// 			"results": [
// 				{"name": "location1", "url": "url1"},
// 				{"name": "location2", "url": "url2"}
// 			]
// 		}`))

// 		loc := GetLocations()
// 		if loc == nil || len(loc.Locations) != 2 {
// 			t.Fatalf("expected 2 locations, got %v", loc)
// 		}
// 	})
// }

// func TestLocationFetching(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	// Mock handler for API responses
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write([]byte(`{
// 			"count": 2,
// 			"next": null,
// 			"previous": null,
// 			"results": [
// 				{"name": "location1", "url": "url1"},
// 				{"name": "location2", "url": "url2"}
// 			]
// 		}`))
// 	}
// 	server := httptest.NewServer(http.HandlerFunc(handler))
// 	defer server.Close()

// 	baseURL = server.URL + "/"

// 	t.Run("Fetching Next and Previous", func(t *testing.T) {
// 		loc := initLocations()
// 		err := loc.GetNext()
// 		if err == nil {
// 			t.Fatalf("expected an error when there is no next page, got %v", loc)
// 		}

// 		err = loc.GetPrevious()
// 		if err == nil {
// 			t.Fatalf("expected an error when there is no previous page, got %v", loc)
// 		}
// 	})

// 	t.Run("Print Locations", func(t *testing.T) {
// 		loc := initLocations()
// 		if err := loc.PrintLocations(); err != nil {
// 			t.Fatalf("expected to print locations without error, but got %v", err)
// 		}
// 	})
// }

// func TestLocationAreaList(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write([]byte(`{
// 			"id": 1,
// 			"name": "test-location",
// 			"areas": [{"name": "area1", "url": "url1"}]
// 		}`))
// 	}
// 	server := httptest.NewServer(http.HandlerFunc(handler))
// 	defer server.Close()

// 	baseURL = server.URL + "/"

// 	t.Run("Get Location Detail", func(t *testing.T) {
// 		locationDetail := GetLocationDetail()
// 		err := locationDetail.GetLocationDetail("test-location")
// 		if err != nil {
// 			t.Fatalf("expected to get location detail without error, but got %v", err)
// 		}

// 		if locationDetail == nil || locationDetail.Name != "test-location" {
// 			t.Fatalf("expected location name to be 'test-location', but got %v", locationDetail.Name)
// 		}

// 	})

// 	t.Run("Location Detail Cached", func(t *testing.T) {
// 		// Adding mocked response to the cache
// 		cache.Add(baseURL+"location-area/test-location", []byte(`{
// 			"id": 1,
// 			"name": "test-location",
// 			"areas": [{"name": "area1", "url": "url1"}]
// 		}`))

// 		locationDetail := GetLocationDetail()
// 		err := locationDetail.GetLocationDetail("test-location")
// 		if err != nil {
// 			t.Fatalf("expected to get location detail from cache without error, but got %v", err)
// 		}

// 		if locationDetail == nil || locationDetail.Name != "test-location" {
// 			t.Fatalf("expected location name to be 'test-location', but got %v", locationDetail.Name)
// 		}
// 	})
// }

// // Utility function to help with cache initialization in tests
// func initCache() *pokecache.Cache {
// 	return pokecache.NewCache(time.Minute * 5)
// }

// func TestFetchFromAPI(t *testing.T) {
// 	cache = initCache()
// 	defer func() { cache = nil }()

// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write([]byte(`{
// 			"id": 1,
// 			"name": "fetch-test"
// 		}`))
// 	}
// 	server := httptest.NewServer(http.HandlerFunc(handler))
// 	defer server.Close()

// 	url := server.URL + "/location-area/fetch-test"
// 	data := new(LocationAreaList)

// 	err := fetchFromAPI(url, data)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if data.Name != "fetch-test" {
// 		t.Fatalf("expected location name 'fetch-test', got %v", data.Name)
// 	}

// 	// Attempt to fetch again and ensure it hits the cache
// 	cachedData := new(LocationAreaList)
// 	err = fetchFromAPI(url, cachedData)
// 	if err != nil {
// 		t.Fatalf("expected no error from cached data, got %v", err)
// 	}

// 	if cachedData.Name != "fetch-test" {
// 		t.Fatalf("expected location name 'fetch-test' from cache, got %v", cachedData.Name)
// 	}
// }

// // Test Unmarshal error case
// func TestUnmarshalError(t *testing.T) {
// 	cache = initCache()
// 	defer func() { cache = nil }()

// 	// Adding data that doesn't match expected structure to cache
// 	cache.Add("bad-data-url", []byte(`{"invalid": "data"}`))

// 	v := new(LocationAreaList)
// 	err := fetchFromAPI("bad-data-url", v)
// 	if err == nil {
// 		t.Fatalf("expected unmarshal error, got nil")
// 	}
// }
