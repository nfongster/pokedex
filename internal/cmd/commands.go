package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CommandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(c *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for name, cmd := range GetRegistry() {
		fmt.Printf("%s: %s\n", name, cmd.Description)
	}
	return nil
}

func CommandMapNext(c *Config) error {
	if c.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	return commandMap(c, c.Next)
}

func CommandMapPrevious(c *Config) error {
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
