package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/nfongster/pokedex/internal/cmd"
)

func main() {
	reader := bufio.NewReader(io.Reader(os.Stdin))
	scanner := bufio.NewScanner(reader)
	config := cmd.Config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
	}

	for {
		fmt.Print("Pokedex > ")
		success := scanner.Scan()
		if !success {
			panic(fmt.Sprintf("Error: %v", scanner.Err()))
		}
		cleanInput := cmd.CleanInput(scanner.Text())

		if len(cleanInput) == 0 {
			fmt.Println("Please type a command.")
		} else if command, exists := cmd.GetRegistry()[cleanInput[0]]; !exists {
			fmt.Println("Unknown command")
		} else if err := command.Callback(&config); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}
