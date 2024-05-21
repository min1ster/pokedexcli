package main

import (
	"bufio"
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

// type locationsPayload struct {
// 	count int
// 	next *string
// 	previous *string
// 	results []location
// }

// type location struct {
// 	name string
// 	url string
// }

func getCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)

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
			currentPage := 0
			err := getLocations(currentPage)
			currentPage += 1
			return err
		},
	}

	commands["help"] = helpCommand
	commands["exit"] = exitCommand
	commands["map"] = mapCommand

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
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", body)
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
