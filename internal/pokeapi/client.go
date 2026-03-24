package pokeapi

import (
	"emagjby/boot-pokedex/internal/pokecache"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const LocationAreaURL = "https://pokeapi.co/api/v2/location-area"
const PokemonURL = "https://pokeapi.co/api/v2/pokemon"

const cacheInterval = 5 * time.Second

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonDetails struct {
	BaseExperience int                `json:"base_experience"`
	Height         int                `json:"height"`
	Name           string             `json:"name"`
	Stats          []PokemonStatEntry `json:"stats"`
	Types          []PokemonTypeEntry `json:"types"`
	Weight         int                `json:"weight"`
}

type PokemonStat struct {
	Name string `json:"name"`
}

type PokemonStatEntry struct {
	BaseStat int         `json:"base_stat"`
	Stat     PokemonStat `json:"stat"`
}

type PokemonType struct {
	Name string `json:"name"`
}

type PokemonTypeEntry struct {
	Type PokemonType `json:"type"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type LocationAreaDetail struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type MapData struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Client struct {
	httpClient *http.Client
	cache      *pokecache.Cache
}

func NewClient() *Client {
	return NewClientWithHTTPClient(&http.Client{})
}

func NewClientWithHTTPClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		cache:      pokecache.NewCache(cacheInterval),
	}
}

func (c *Client) GetLocationAreas(url string) (MapData, error) {
	body, err := c.get(url)
	if err != nil {
		return MapData{}, err
	}

	var out MapData
	if err := json.Unmarshal(body, &out); err != nil {
		return MapData{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return out, nil
}

func (c *Client) GetLocationArea(url string) (LocationAreaDetail, error) {
	body, err := c.get(fmt.Sprintf("%s/%s", LocationAreaURL, url))
	if err != nil {
		return LocationAreaDetail{}, err
	}

	var out LocationAreaDetail
	if err := json.Unmarshal(body, &out); err != nil {
		return LocationAreaDetail{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return out, nil
}

func (c *Client) GetPokemon(name string) (PokemonDetails, error) {
	body, err := c.get(fmt.Sprintf("%s/%s", PokemonURL, name))
	if err != nil {
		return PokemonDetails{}, err
	}

	var out PokemonDetails
	if err := json.Unmarshal(body, &out); err != nil {
		return PokemonDetails{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return out, nil
}

func (c *Client) get(url string) ([]byte, error) {
	if cachedBody, ok := c.cache.Get(url); ok {
		return cachedBody, nil
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.cache.Add(url, body)

	return body, nil
}
