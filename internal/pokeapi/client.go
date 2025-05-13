package pokeapi

import (
	"encoding/json"
	_ "errors"
	_ "io" // don't need?
	"log"
	"net/http"
)

type Config struct {
	NextURL    string // no commas in structs in go! only semicolons if all in one line
	PrevURL    string
	PokeClient *Client
}

type LocationAreaPage struct { // for holding data for return from ListLocationAreas()
	Names   []string // slice of area names
	NextURL string   // next and prev for assigning to config after return.
	PrevURL string
}

type PokeAPI struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []locationSummary `json:"results"`
}

type locationSummary struct {
	Name string `json:"name"`
}

type Client struct {
	httpClient http.Client
}

func NewClient() *Client {
	return &Client{
		//initialize data
		httpClient: http.Client{},
	}
}

func (c *Client) ListLocationAreas(url string) (LocationAreaPage, error) {
	//if url empty, use base
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	returnData := LocationAreaPage{}

	// Make HTTP request to the URL

	res, err := http.Get(url)
	if err != nil {
		return LocationAreaPage{}, err
	}
	/*
		//set req headers?
		body, err := io.ReadAll(res.Body)
	*/

	defer res.Body.Close() // do i need to defer?

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d", res.StatusCode)
	}

	if err != nil {
		return LocationAreaPage{}, err
	}

	// Parse the response

	var pokeapi PokeAPI // LocationAreaPage?
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokeapi) // not sure if correct
	if err != nil {
		return LocationAreaPage{}, err
	}

	// Return the data and update Config // actuall, DON'T update config, but returning values for caller to update.
	returnData.PrevURL = pokeapi.Previous
	returnData.NextURL = pokeapi.Next

	for _, loc := range pokeapi.Results {
		returnData.Names = append(returnData.Names, loc.Name)
	}

	// update config...?
	return returnData, nil

}
