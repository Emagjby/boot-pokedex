package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type mapData struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []locationArea `json:"results"`
}

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func(*config) error
}

type config struct {
	next     string
	previous string
}

var commands map[string]cliCommand

func initCommands(config *config) {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show this help message",
			config:      config,
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show the map of the Pokemon world",
			config:      config,
			callback:    commandMap,
		},
		"mapb": {
			name:		 "mapb",
			description: "Show the previous page of the map",
			config:		 config,
			callback:	 commandMapb,
		},
	}
}

func initConfig() config {
	return config{
		next:     "",
		previous: "",
	}
}

func main() {
	config := initConfig()
	initCommands(&config)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Internal Error:", err)
			continue
		}

		err = handleInput(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func handleInput(input string) error {
	parsed := cleanInput(input)
	if len(parsed) == 0 {
		return nil
	}

	cmd, ok := commands[parsed[0]]
	if !ok {
		return fmt.Errorf("Unknown command: %q", parsed[0])
	}

	return cmd.callback(cmd.config)
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	var body []byte
	var err error
	if cfg.next == "" {
		body, err = getFromApi("https://pokeapi.co/api/v2/location-area")
	} else {
		body, err = getFromApi(cfg.next)
	}
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	data, err := formatMapData(body)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	cfg.next = data.Next
	cfg.previous = data.Previous
	return nil
}

func commandMapb(cfg *config) error {
	var body []byte
	var err error
	if cfg.previous == "" {
		body, err = getFromApi("https://pokeapi.co/api/v2/location-area")
	} else {
		body, err = getFromApi(cfg.previous)
	}
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	data, err := formatMapData(body)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	cfg.next = data.Next
	cfg.previous = data.Previous
	return nil
}

func formatMapData(data []byte) (mapData, error) {
	var out mapData
	err := json.Unmarshal(data, &out)
	if err != nil {
		return mapData{}, fmt.Errorf("Failed to decode map data: %v", err)
	}

	return out, nil
}

func getFromApi(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get map data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get map data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	return body, nil
}
