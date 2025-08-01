package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/nfongster/pokedex/internal/pokecache"
)

func CommandExit(c *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(c *Config, args ...string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for name, cmd := range GetRegistry() {
		fmt.Printf("%s: %s\n", name, cmd.Description)
	}
	return nil
}

func CommandMapNext(c *Config, args ...string) error {
	if c.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	return commandMap(c, c.Next)
}

func CommandMapPrevious(c *Config, args ...string) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return commandMap(c, c.Previous)
}

func CommandExplore(c *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a location to explore")
	}

	location := args[0]
	url := getLocationAreaUrl(location)
	bytes, err := getByteArray(&c.Cache, url)
	if err != nil {
		return err
	}
	if _, cached := c.Cache.Get(url); !cached {
		//fmt.Printf("\033[32mAdding data from %s to cache\n\033[0m", url)
		c.Cache.Add(url, bytes)
	}

	areaData := LocationArea{}
	if err := json.Unmarshal(bytes, &areaData); err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\nFound Pokemon:\n", location)
	for _, encounter := range areaData.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func CommandCatch(c *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a pokemon to catch")
	}

	name := args[0]
	url := getPokemonUrl(name)
	bytes, err := getByteArray(&c.Cache, url)
	if err != nil {
		return err
	}
	if _, cached := c.Cache.Get(url); !cached {
		//fmt.Printf("\033[32mAdding data from %s to cache\n\033[0m", url)
		c.Cache.Add(url, bytes)
	}

	pokemon := Pokemon{}
	if err := json.Unmarshal(bytes, &pokemon); err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	if caught := tryCatchPokemon(pokemon.BaseExperience); caught {
		fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", name)
		c.Pokedex[name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}

func tryCatchPokemon(exp int) bool {
	//fmt.Printf(" Base Exp: %d\n", exp)
	// Scale base experience (typically 50-300) to a catch rate
	maxChance := 0.8 // 80% for easiest Pokemon
	minChance := 0.1 // 10% for hardest Pokemon
	catchChance := maxChance - (float64(exp-50)/250.0)*(maxChance-minChance)
	return rand.Float64() < catchChance
}

func CommandInspect(c *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a pokemon to catch")
	}

	name := args[0]
	pokemon, exists := c.Pokedex[name]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("  - %s\n", typ.Type.Name)
	}
	return nil
}

func CommandPokedex(c *Config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for name := range c.Pokedex {
		fmt.Printf("  - %s\n", name)
	}
	return nil
}

func commandMap(c *Config, url string) error {
	bytes, err := getByteArray(&c.Cache, url)
	if err != nil {
		return err
	}
	if _, cached := c.Cache.Get(url); !cached {
		//fmt.Printf("\033[32mAdding data from %s to cache\n\033[0m", url)
		c.Cache.Add(url, bytes)
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

func getByteArray(c *pokecache.Cache, url string) ([]byte, error) {
	if data, cached := c.Get(url); cached {
		//fmt.Println("\033[32mCache hit\033[0m")
		return data, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	return io.ReadAll(res.Body)
}
