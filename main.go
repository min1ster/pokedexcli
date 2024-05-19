package main

import (
	"fmt"
	"bufio"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

func getCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)

	helpCommand := cliCommand{
		name: "help",
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
		name: "exit",
		description: "Exit the Pokedex",
		callback: func() error {
			os.Exit(0)
			return nil
		},
	}

	commands["help"] = helpCommand
	commands["exit"] = exitCommand

	return commands
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