package locations

import (
	"fmt"
	"net/http"
	"log"
	"io"
	"encoding/json"
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

func GetLocations(page int) error {
	offset := 20 * page
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", offset)
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
	locations := locationsPayload{}
	json.Unmarshal(bodyBytes, &locations)
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}