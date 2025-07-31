package cmd

import "strings"

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
			Description: "Explore a given area.  Pass the area as an argument.",
			Callback:    CommandExplore,
		},
	}
}

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
