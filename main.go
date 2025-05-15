package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gainax2k1/pokedexcli/internal/pokeapi"
	_ "github.com/gainax2k1/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *pokeapi.Config) error
}

// Declare the variable without initializing it
var commandMap map[string]cliCommand

// Init function runs after variable declaration but before main
func init() {
	commandMap = map[string]cliCommand{
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
			description: "Displays names of locations, going forward a page", //my descript
			callback:    commandListMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays names of locations, going back a page", // my descript
			callback:    commandListMapBack,
		},
	}
}

/* ******************************************
When you come back, you'll be focusing on:

   1. cache
 ****************************************** */

func main() {
	//old for initial testing: fmt.Println("Hello, World!")

	cfg := &pokeapi.Config{
		PokeClient: pokeapi.NewClient(5 * time.Minute),
	}

	userInputScanner := bufio.NewScanner(os.Stdin) // correct?

	for {
		fmt.Print("Pokedex > ") //prints w/o using a newline
		userInputScanner.Scan()
		userInput := userInputScanner.Text()

		cleanedInput := cleanInput(userInput)

		if len(cleanedInput) > 0 {

			command := cleanedInput[0]

			if cmd, exists := commandMap[command]; exists {
				err := cmd.callback(cfg)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func cleanInput(text string) []string {
	//creates "lowered", trimming  trailing whitespace and making lowercase
	lowered := strings.ToLower(strings.TrimSpace(text))

	words := strings.Fields(lowered)

	return words
}

func commandExit(cfg *pokeapi.Config) error { // for immmediately quitting the program
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *pokeapi.Config) error { //for printing help text
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cliVal := range commandMap {
		fmt.Printf("%s: %s\n", cliVal.name, cliVal.description)

	}

	return nil
}

func commandListMap(cfg *pokeapi.Config) error {

	// Call the API using the client
	locationAreaPage, err := cfg.PokeClient.ListLocationAreas(cfg.NextURL)
	if err != nil {
		return err
	}

	// Update the config with the new URLs
	cfg.NextURL = locationAreaPage.NextURL
	cfg.PrevURL = locationAreaPage.PrevURL

	// Print the location names
	for _, location := range locationAreaPage.Names {
		fmt.Println(location)
	}

	return nil
}

func commandListMapBack(cfg *pokeapi.Config) error {
	if cfg.PrevURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	// Call the API using the client
	locationAreaPage, err := cfg.PokeClient.ListLocationAreas(cfg.PrevURL)
	if err != nil {
		return err
	}

	// Update the config with the new URLs
	cfg.NextURL = locationAreaPage.NextURL
	cfg.PrevURL = locationAreaPage.PrevURL

	// Print the location names
	for _, location := range locationAreaPage.Names {
		fmt.Println(location)
	}

	return nil

}
