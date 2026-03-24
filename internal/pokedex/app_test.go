package pokedex

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"emagjby/boot-pokedex/internal/pokeapi"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestHandleInputExplore(t *testing.T) {
	var output bytes.Buffer
	app := New(strings.NewReader(""), &output)
	app.client = pokeapi.NewClientWithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"name":"pastoria-city-area","pokemon_encounters":[{"pokemon":{"name":"wooper","url":"https://pokeapi.co/api/v2/pokemon/194/"}},{"pokemon":{"name":"quagsire","url":"https://pokeapi.co/api/v2/pokemon/195/"}}]}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		}),
	})

	if err := app.handleInput("explore pastoria-city-area\n"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Exploring pastoria-city-area...\nFound Pokemon:\n - wooper\n - quagsire\n"
	if output.String() != expected {
		t.Fatalf("expected output %q, got %q", expected, output.String())
	}
}

func TestHandleInputExploreRequiresAreaName(t *testing.T) {
	app := New(strings.NewReader(""), io.Discard)

	err := app.handleInput("explore\n")
	if err == nil {
		t.Fatal("expected an error for missing area name")
	}

	if err.Error() != "usage: explore <location-area-name>" {
		t.Fatalf("expected usage error, got %q", err.Error())
	}
}

func TestHandleInputCatchSuccess(t *testing.T) {
	var output bytes.Buffer
	app := New(strings.NewReader(""), &output)
	app.client = pokeapi.NewClientWithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"name":"pikachu","base_experience":112}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		}),
	})
	app.catchRoll = func(_ int) int { return 0 }

	if err := app.handleInput("catch pikachu\n"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Throwing a Pokeball at pikachu...\npikachu was caught!\n"
	if output.String() != expected {
		t.Fatalf("expected output %q, got %q", expected, output.String())
	}

	if _, ok := app.pokedex["pikachu"]; !ok {
		t.Fatal("expected pikachu to be added to the pokedex")
	}
}

func TestHandleInputCatchEscape(t *testing.T) {
	var output bytes.Buffer
	app := New(strings.NewReader(""), &output)
	app.client = pokeapi.NewClientWithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"name":"pikachu","base_experience":112}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		}),
	})
	app.catchRoll = func(_ int) int { return 99 }

	if err := app.handleInput("catch pikachu\n"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Throwing a Pokeball at pikachu...\npikachu escaped!\n"
	if output.String() != expected {
		t.Fatalf("expected output %q, got %q", expected, output.String())
	}

	if _, ok := app.pokedex["pikachu"]; ok {
		t.Fatal("did not expect escaped pokemon to be added to the pokedex")
	}
}

func TestHandleInputCatchRequiresPokemonName(t *testing.T) {
	app := New(strings.NewReader(""), io.Discard)

	err := app.handleInput("catch\n")
	if err == nil {
		t.Fatal("expected an error for missing pokemon name")
	}

	if err.Error() != "usage: catch <pokemon-name>" {
		t.Fatalf("expected usage error, got %q", err.Error())
	}
}
