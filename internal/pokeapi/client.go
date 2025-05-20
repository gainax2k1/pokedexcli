package pokeapi

import (
	"encoding/json"
	_ "errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gainax2k1/pokedexcli/internal/pokecache"
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

type PokeEncounters struct { //top level
	Encounters []PokeList `json:"pokemon_encounters"`
}
type PokeList struct { // mid level
	Pokemon PokemonName `json:"pokemon"`
}

type PokemonName struct { //base level
	Name string `json:"name"` //// resume here! ************
}

type PokemonStat struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokeStats struct {
	Base_Stat int         `json:"base_stat"`
	Effort    int         `json:"effort"`
	Stat      PokemonStat `json:"stat"`
}

type PokemonType struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokeTypes struct {
	Slot int         `json:"slot"`
	Type PokemonType `json:"type"`
}
type PokemonInfo struct { // extra info commented out for now, but maybe neccessary in the future?
	Id              int         `json:"id"`
	Name            string      `json:"name"`
	Base_Experience int         `json:"base_experience"`
	Height          int         `json:"height"`
	Is_Default      bool        `json:"is_default"`
	Order           int         `json:"order"`
	Weight          int         `json:"weight"`
	Stats           []PokeStats `json:"stats"`
	Types           []PokeTypes `json:"types"`
}

type Client struct {
	cache      *pokecache.Cache // for time caching
	httpClient http.Client
}

func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		//initialize data

		cache:      pokecache.NewCache(cacheInterval),
		httpClient: http.Client{},
	}
}

func (c *Client) ListLocationAreas(url string) (LocationAreaPage, error) {
	//if url empty, use base
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	returnData := LocationAreaPage{}
	// First check if we have this URL in the cache
	var body []byte
	var err error

	cachedData, found := c.cache.Get(url)
	if found {
		// Use the cached data
		fmt.Println("Cache hit!") // Optional logging
		body = cachedData
	} else {
		// Cache miss - make the HTTP request
		fmt.Println("Cache miss! Fetching from API...") // Optional logging

		res, err := http.Get(url)
		if err != nil {
			return LocationAreaPage{}, err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return LocationAreaPage{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

		// Read the response body
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationAreaPage{}, err
		}

		// Add to cache
		c.cache.Add(url, body)
	}

	// Parse the response

	var pokeapi PokeAPI

	err = json.Unmarshal(body, &pokeapi)
	if err != nil {
		return LocationAreaPage{}, err
	}

	// Return the data and update Config // actuall, DON'T update config, but returning values for caller to update.
	returnData.PrevURL = pokeapi.Previous
	returnData.NextURL = pokeapi.Next
	for _, loc := range pokeapi.Results {
		returnData.Names = append(returnData.Names, loc.Name)
	}

	return returnData, nil

}

func (c *Client) ListPokemon(arg []string) ([]PokemonName, error) {
	pokeList := []PokemonName{}

	// insert stuff here!

	url := "https://pokeapi.co/api/v2/location-area/" + arg[0] //append city name to url

	// First check if we have this URL in the cache
	var body []byte
	var err error

	cachedData, found := c.cache.Get(url)
	if found {
		// Use the cached data
		fmt.Println("Cache hit!") // Optional logging
		body = cachedData
	} else {
		// Cache miss - make the HTTP request
		fmt.Println("Cache miss! Fetching from API...") // Optional logging

		res, err := http.Get(url)
		if err != nil {
			return []PokemonName{}, err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return []PokemonName{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

		// Read the response body
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return []PokemonName{}, err
		}

		// Add to cache
		c.cache.Add(url, body)
	}

	// Parse the response

	var enounters PokeEncounters

	err = json.Unmarshal(body, &enounters)
	if err != nil {
		return []PokemonName{}, err
	}

	for _, pList := range enounters.Encounters {
		pokeList = append(pokeList, pList.Pokemon)
	}

	// update config...?
	return pokeList, nil

}

func (c *Client) CatchPokemon(arg []string) (bool, PokemonInfo, error) {

	successful := false // successful = true if pokemon caught successfully, default to false for err conditions

	url := "https://pokeapi.co/api/v2/pokemon/" + arg[0] //append pokemon name to url

	// First check if we have this URL in the cache
	var body []byte
	var err error

	cachedData, found := c.cache.Get(url)
	if found {
		// Use the cached data
		fmt.Println("Cache hit!") // Optional logging
		body = cachedData
	} else {
		// Cache miss - make the HTTP request
		fmt.Println("Cache miss! Fetching from API...") // Optional logging

		res, err := http.Get(url)
		if err != nil {
			return successful, PokemonInfo{}, err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return successful, PokemonInfo{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

		// Read the response body
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return successful, PokemonInfo{}, err
		}

		// Add to cache
		c.cache.Add(url, body)
	}

	// Parse the response

	var pokeInfo PokemonInfo

	err = json.Unmarshal(body, &pokeInfo)
	if err != nil {
		return successful, PokemonInfo{}, err
	}

	value := 10000.0 * (1.0 / float64(pokeInfo.Base_Experience)) * rand.Float64() // catching algorithm

	fmt.Printf("value: %f\n", value)
	if value > 50 {
		successful = true
	}

	return successful, pokeInfo, nil

}
