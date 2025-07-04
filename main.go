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
	callback    func(cfg *pokeapi.Config, arg []string) error
}

// Declare the variable without initializing it
var commandMap map[string]cliCommand

/*
type Pokemon struct {
	name string
}*/

var pokedex = map[string]pokeapi.PokemonInfo{} // to store all the users caught pokemon

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
		"explore": {
			name:        "explore",
			description: "Displays list of pokeman from the specified location", // my descript
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch specified pokemon", // my descript
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a specified pokemon if already caught", // my descript
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught pokemon", // my descript
			callback:    commandPokedex,
		},
	}
}

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
			argument := []string{}
			if len(cleanedInput) > 1 {
				argument = cleanedInput[1:]
			}

			if cmd, exists := commandMap[command]; exists {
				err := cmd.callback(cfg, argument)
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

func commandExit(cfg *pokeapi.Config, _ []string) error { // for immmediately quitting the program
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *pokeapi.Config, _ []string) error { //for printing help text
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cliVal := range commandMap {
		fmt.Printf("%s: %s\n", cliVal.name, cliVal.description)

	}

	return nil
}

func commandListMap(cfg *pokeapi.Config, _ []string) error {

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

func commandListMapBack(cfg *pokeapi.Config, _ []string) error {
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

func commandExplore(cfg *pokeapi.Config, arg []string) error {
	// in case no city is specified for exploring
	if len(arg) == 0 {
		fmt.Println("No city specified for exploring.")
		fmt.Println("Command should be in format 'explore city-name'")
		return nil
	}
	fmt.Printf("Exploring %s...\n", arg[0])
	fmt.Println("Found Pokemon:")
	// Call the API using the client
	pokeList, err := cfg.PokeClient.ListPokemon(arg)
	if err != nil {
		return err
	}

	for _, pokemon := range pokeList {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil

}

func commandCatch(cfg *pokeapi.Config, arg []string) error {
	// in case no pokemons specified for attempt to catch
	if len(arg) == 0 {
		fmt.Println("No pokemon specified for catching.")
		fmt.Println("Command should be in format 'catch pokemon-name'")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", arg[0])

	// insert catching logic here
	succesfulCatch, pokeInfo, err := cfg.PokeClient.CatchPokemon(arg)
	if err != nil {
		return err
	}

	if succesfulCatch {
		fmt.Printf("%s was caught!\n", arg[0])
		fmt.Println("You may now inspect it with the inspect command.")
		pokedex[arg[0]] = pokeInfo

	} else {
		fmt.Printf("%s escaped!\n", arg[0])
	}
	return nil
}

func commandInspect(cfg *pokeapi.Config, arg []string) error {
	// in case no pokemons specified for inspection
	if len(arg) == 0 {
		fmt.Println("No pokemon specified for inspection.")
		fmt.Println("Command should be in format 'inspect pokemon-name'")
		return nil
	}
	pokemonInfo, success := pokedex[arg[0]]
	if !success {
		fmt.Println("you have not yet caught that pokemon")
		return nil
	} else {
		printInfo(pokemonInfo)
	}

	return nil
}

func printInfo(pokemonInfo pokeapi.PokemonInfo) {
	fmt.Printf("Name: %s\n", pokemonInfo.Name)
	fmt.Printf("Height: %d\n", pokemonInfo.Height)
	fmt.Printf("Weight: %d\n", pokemonInfo.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonInfo.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.Base_Stat)
	}

	fmt.Println("Types:")

	for _, ptype := range pokemonInfo.Types {
		fmt.Printf("  - %s\n", ptype.Type.Name)
	}
}

func commandPokedex(cfg *pokeapi.Config, arg []string) error {
	fmt.Println("Your Pokedex:")
	for k := range pokedex {
		fmt.Printf(" - %s\n", k)
	}
	return nil
}
