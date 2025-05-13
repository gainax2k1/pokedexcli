package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func main() {
	//old for initial testing: fmt.Println("Hello, World!")

	userInputScanner := bufio.NewScanner(os.Stdin) // correct?

	for {
		fmt.Print("Pokedex > ") //prints w/o using a newline
		userInputScanner.Scan()
		userInput := userInputScanner.Text()

		cleanedInput := cleanInput(userInput)

		if len(cleanedInput) > 0 {

			command := cleanedInput[0]

			if cmd, exists := commandMap[command]; exists {
				err := cmd.callback()
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
		/* for testing input
		        if len(cleanInput) > 0 {
					fmt.Printf("Your command was: %s\n", cleanInput[0])
				}
		*/
	}
}

func cleanInput(text string) []string {
	//creates "lowered", trimming  trailing whitespace and making lowercase
	lowered := strings.ToLower(strings.TrimSpace(text))

	words := strings.Fields(lowered)

	return words
}

func commandExit() error { // for immmediately quitting the program
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error { //for printing help text
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cliVal := range commandMap {
		fmt.Printf("%s: %s\n", cliVal.name, cliVal.description)

	}

	return nil
}
