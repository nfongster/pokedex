package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
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
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapPrevious,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for name, cmd := range getRegistry() {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

type LocationBatch struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func commandMapNext(c *Config) error {
	if c.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	return commandMap(c, c.Next)
}

func commandMapPrevious(c *Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return commandMap(c, c.Previous)
}

func commandMap(c *Config, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d", res.StatusCode)
	} else if err != nil {
		return err
	}

	batch := LocationBatch{}
	if err := json.Unmarshal(bytes, &batch); err != nil {
		return err
	}

	for _, location := range batch.Results {
		fmt.Println(location.Name)
	}

	c.Next = batch.Next
	c.Previous = batch.Previous
	return nil
}

func main() {
	reader := bufio.NewReader(io.Reader(os.Stdin))
	scanner := bufio.NewScanner(reader)
	config := Config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
	}

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
		} else if err := command.callback(&config); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}
