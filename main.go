package main

import (
	"fmt"
	"strings"
)

func main() {

	fmt.Println("Hello, World!")

}

func cleanInput(text string) []string {
	//creates "lowered", trimming  trailing whitespace and making lowercase
	lowered := strings.ToLower(strings.TrimSpace(text))

	words := strings.Fields(lowered)

	return words // placeholder
}
