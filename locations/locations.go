package locations

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/min1ster/pokedexcli/pokecache"
)

type locationsPayload struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func GetLocations(page int, cache *pokecache.Cache) error {
	offset := 20 * page
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", offset)
	cacheEntry, ok := cache.Entries[endpoint]
	if ok {
		cache.Add(endpoint, cacheEntry.Val)
		return handleOutput(cacheEntry.Val)
	}

	res, err := http.Get(endpoint)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Response failed with status code: %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	cache.Add(endpoint, bodyBytes)
	return handleOutput(bodyBytes)
}

func handleOutput(bodyBytes []byte) error {
	locations := locationsPayload{}
	json.Unmarshal(bodyBytes, &locations)
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}
