package locations

import (
	"encoding/json"
	"fmt"
	"github.com/min1ster/pokedexcli/pokecache"
	"io"
	"log"
	"net/http"
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

type locationPayload struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Callback func(bodyBytes []byte) error

func GetLocations(page int, cache *pokecache.Cache) error {
	offset := 20 * page
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", offset)
	return endointHandler(cache, endpoint, handleLocationsOutput)
}

func GetLocation(name string, cache *pokecache.Cache) error {
	fmt.Printf("Exploring %s...\n", name)
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
	return endointHandler(cache, endpoint, handleLocationOutput)
}

func endointHandler(cache *pokecache.Cache, endpoint string, callback Callback) error {
	cacheEntry, ok := cache.Entries[endpoint]
	if ok {
		cache.Add(endpoint, cacheEntry.Val)
		return callback(cacheEntry.Val)
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
	return callback(bodyBytes)
}

func handleLocationsOutput(bodyBytes []byte) error {
	locations := locationsPayload{}
	json.Unmarshal(bodyBytes, &locations)
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func handleLocationOutput(bodyBytes []byte) error {
	pokemon := locationPayload{}
	json.Unmarshal(bodyBytes, &pokemon)
	fmt.Println("Found Pokemon:")
	for _, record := range pokemon.PokemonEncounters {
		fmt.Printf(" - %s\n", record.Pokemon.Name)
	}
	return nil
}
