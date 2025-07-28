package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
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
		firstWord := ""
		if len(cleanInput) > 0 {
			firstWord = cleanInput[0]
		}
		fmt.Printf("Your command was: %s\n", firstWord)
	}
}
