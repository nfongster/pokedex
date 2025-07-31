package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nfongster/pokedex/internal/cmd"
	"github.com/nfongster/pokedex/internal/pokecache"
)

func main() {
	reader := bufio.NewReader(io.Reader(os.Stdin))
	scanner := bufio.NewScanner(reader)
	config := cmd.Config{
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		Previous: "",
		Pokedex:  map[string]cmd.Pokemon{},
		Cache:    *pokecache.NewCache(5 * time.Second),
	}

	for {
		fmt.Print("Pokedex > ")
		success := scanner.Scan()
		if !success {
			panic(fmt.Sprintf("Error: %v", scanner.Err()))
		}
		cleanInput := cmd.CleanInput(scanner.Text())
		args := []string{}
		if len(cleanInput) > 1 {
			args = cleanInput[1:]
		}

		if len(cleanInput) == 0 {
			fmt.Println("Please type a command.")
		} else if command, exists := cmd.GetRegistry()[cleanInput[0]]; !exists {
			fmt.Println("Unknown command")
		} else if err := command.Callback(&config, args...); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}
