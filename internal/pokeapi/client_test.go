package pokeapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gainax2k1/pokedexcli/internal/pokeapi"
)

func TestNewClient(t *testing.T) {

	client := pokeapi.NewClient(5 * time.Second)
	if client == nil {
		t.Error("expected client to not be nil")
	}
}
func TestListLocationAreas(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"count": 2,
			"next": "https://example.com/api/v2/location-area?offset=20",
			"previous": null,
			"results": [
				{"name": "test-area-1"},
				{"name": "test-area-2"}
			]
		}`))
	}))
	defer server.Close()

	// Create a client and call the method with our test server URL
	client := pokeapi.NewClient(5 * time.Second)
	result, err := client.ListLocationAreas(server.URL)

	// Check the results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the response was parsed correctly
	if len(result.Names) != 2 {
		t.Errorf("Expected 2 location areas, got %d", len(result.Names))
	}

}
