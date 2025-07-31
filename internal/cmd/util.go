package cmd

import (
	"fmt"
	"strings"
)

func GetRegistry() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Displays the next 20 locations",
			Callback:    CommandMapNext,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous 20 locations",
			Callback:    CommandMapPrevious,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore an area.  Pass the name of the area as an argument.",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a pokemon.  Pass the name of the pokemon as an argument.",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a pokemon's stats.",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Display a list of all the pokemon you've caught.",
			Callback:    CommandPokedex,
		},
	}
}

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func getLocationAreaUrl(location string) string {
	return fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)
}

func getPokemonUrl(name string) string {
	return fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
}
