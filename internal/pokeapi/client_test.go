package pokeapi

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"emagjby/boot-pokedex/internal/pokecache"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestGetLocationArea(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != LocationAreaURL+"/pastoria-city-area" {
					t.Fatalf("expected URL %q, got %q", LocationAreaURL+"/pastoria-city-area", req.URL.String())
				}

				body := `{"name":"pastoria-city-area","pokemon_encounters":[{"pokemon":{"name":"wooper","url":"https://pokeapi.co/api/v2/pokemon/194/"}},{"pokemon":{"name":"quagsire","url":"https://pokeapi.co/api/v2/pokemon/195/"}}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
		cache: pokecache.NewCache(time.Minute),
	}

	area, err := client.GetLocationArea("pastoria-city-area")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if area.Name != "pastoria-city-area" {
		t.Fatalf("expected area name %q, got %q", "pastoria-city-area", area.Name)
	}

	if len(area.PokemonEncounters) != 2 {
		t.Fatalf("expected 2 encounters, got %d", len(area.PokemonEncounters))
	}

	if area.PokemonEncounters[0].Pokemon.Name != "wooper" {
		t.Fatalf("expected first pokemon %q, got %q", "wooper", area.PokemonEncounters[0].Pokemon.Name)
	}
	if area.PokemonEncounters[1].Pokemon.Name != "quagsire" {
		t.Fatalf("expected second pokemon %q, got %q", "quagsire", area.PokemonEncounters[1].Pokemon.Name)
	}
}

func TestGetPokemon(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != PokemonURL+"/pikachu" {
					t.Fatalf("expected URL %q, got %q", PokemonURL+"/pikachu", req.URL.String())
				}

				body := `{"name":"pikachu","base_experience":112}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
		cache: pokecache.NewCache(time.Minute),
	}

	pokemon, err := client.GetPokemon("pikachu")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if pokemon.Name != "pikachu" {
		t.Fatalf("expected pokemon name %q, got %q", "pikachu", pokemon.Name)
	}

	if pokemon.BaseExperience != 112 {
		t.Fatalf("expected base experience %d, got %d", 112, pokemon.BaseExperience)
	}
}
