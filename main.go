package main

import (
	"bufio"
	"fmt"
	"github.com/min1ster/pokedexcli/locations"
	"github.com/min1ster/pokedexcli/pokecache"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
}

func getCommands() map[string]cliCommand {
	interval := time.Minute * 2
	cache := pokecache.NewCache(interval)
	commands := make(map[string]cliCommand)
	currentPage := -1

	helpCommand := cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(argument string) error {
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
		callback: func(argument string) error {
			os.Exit(0)
			return nil
		},
	}

	mapCommand := cliCommand{
		name:        "map",
		description: "Displays the names of 20 location areas in the Pokemon world. Each call displays the next 20 locations.",
		callback: func(argument string) error {
			currentPage += 1
			err := locations.GetLocations(currentPage, cache)
			return err
		},
	}

	mapBCommand := cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 location areas in the Pokemon world. Each call displays the previous 20 locations.",
		callback: func(argument string) error {
			if currentPage > -1 {
				currentPage -= 1
			}
			err := locations.GetLocations(currentPage, cache)
			return err
		},
	}

	exploreCommand := cliCommand{
		name:        "explore",
		description: "Displays the pokemon available at a given location.",
		callback: func(location string) error {
			err := locations.GetLocation(location, cache)
			return err
		},
	}

	commands["help"] = helpCommand
	commands["exit"] = exitCommand
	commands["map"] = mapCommand
	commands["mapb"] = mapBCommand
	commands["explore"] = exploreCommand

	return commands
}

func main() {
	commands := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()
			args := strings.Fields(input)
			commandName := args[0]
			var commandArgument string
			if len(args) > 1 {
				commandArgument = args[1]
			}
			command, ok := commands[commandName]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}
			command.callback(commandArgument)
		}
	}
}
