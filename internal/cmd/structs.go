package cmd

import "github.com/nfongster/pokedex/internal/pokecache"

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

type LocationBatch struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     string
	Previous string
	Cache    pokecache.Cache
}
