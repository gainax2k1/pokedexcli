package pokeapi

import (
	"encoding/json"
	_ "errors"
	"fmt"
	"io"
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

	/* before caching implemented
	// Make HTTP request to the URL

	res, err := http.Get(url)
	if err != nil {
		return LocationAreaPage{}, err
	}

	defer res.Body.Close() // do i need to defer?

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d", res.StatusCode)
	}

	if err != nil {
		return LocationAreaPage{}, err
	}
	*/

	// Parse the response

	var pokeapi PokeAPI // LocationAreaPage?

	err = json.Unmarshal(body, &pokeapi)
	if err != nil {
		return LocationAreaPage{}, err
	}

	/*

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&pokeapi) // not sure if correct
		if err != nil {
			return LocationAreaPage{}, err
		} */

	// Return the data and update Config // actuall, DON'T update config, but returning values for caller to update.
	returnData.PrevURL = pokeapi.Previous
	returnData.NextURL = pokeapi.Next

	for _, loc := range pokeapi.Results {
		returnData.Names = append(returnData.Names, loc.Name)
	}

	// update config...?
	return returnData, nil

}
