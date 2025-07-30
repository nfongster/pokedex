package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getRegistry() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for name, cmd := range getRegistry() {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

func main() {

	reader := bufio.NewReader(io.Reader(os.Stdin))
	scanner := bufio.NewScanner(reader)
	for {
		fmt.Print("Pokedex > ")
		success := scanner.Scan()
		if !success {
			panic(fmt.Sprintf("Error: %v", scanner.Err()))
		}
		cleanInput := cleanInput(scanner.Text())

		if len(cleanInput) == 0 {
			fmt.Println("Please type a command.")
		} else if command, exists := getRegistry()[cleanInput[0]]; !exists {
			fmt.Println("Unknown command")
		} else if err := command.callback(); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}
