package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type locationsPayload struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func getCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)
	currentPage := -1

	helpCommand := cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func() error {
			fmt.Print("\nWelcome to Pokedex!\n\n")
			for command := range commands {
				fmt.Printf("%s: %s\n", commands[command].name, commands[command].description)
			}
			fmt.Println()
			return nil
		},
	}

	exitCommand := cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback: func() error {
			os.Exit(0)
			return nil
		},
	}

	mapCommand := cliCommand{
		name:        "map",
		description: "Displays the names of 20 location areas in the Pokemon world. Each call displays the next 20 locations.",
		callback: func() error {
			currentPage += 1
			err := getLocations(currentPage)
			return err
		},
	}

	mapBCommand := cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 location areas in the Pokemon world. Each call displays the previous 20 locations.",
		callback: func() error {
			if currentPage > -1 {
				currentPage -= 1
			}
			err := getLocations(currentPage)
			return err
		},
	}

	commands["help"] = helpCommand
	commands["exit"] = exitCommand
	commands["map"] = mapCommand
	commands["mapb"] = mapBCommand

	return commands
}

func getLocations(page int) error {
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

func main() {
	commands := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()
			command, ok := commands[input]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}
			command.callback()
		}
	}
}
