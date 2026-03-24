package pokedex

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"time"

	"emagjby/boot-pokedex/internal/pokeapi"
	"emagjby/boot-pokedex/internal/repl"
)

type command struct {
	name        string
	description string
	callback    func([]string) error
}

type config struct {
	next     string
	previous string
}

type App struct {
	reader    *bufio.Reader
	writer    io.Writer
	client    *pokeapi.Client
	rng       *rand.Rand
	catchRoll func(int) int
	pokedex   map[string]pokeapi.PokemonDetails
	config    config
	commands  map[string]command
}

func New(input io.Reader, output io.Writer) *App {
	app := &App{
		reader:  bufio.NewReader(input),
		writer:  output,
		client:  pokeapi.NewClient(),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		pokedex: map[string]pokeapi.PokemonDetails{},
		config:  config{},
	}
	app.catchRoll = app.rng.Intn

	app.commands = map[string]command{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    app.commandExit,
		},
		"help": {
			name:        "help",
			description: "Show this help message",
			callback:    app.commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show the map of the Pokemon world",
			callback:    app.commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous page of the map",
			callback:    app.commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    app.commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon",
			callback:    app.commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Show details for a caught Pokemon",
			callback:    app.commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List caught Pokemon",
			callback:    app.commandPokedex,
		},
	}

	return app
}

func (a *App) Run() error {
	for {
		fmt.Fprint(a.writer, "Pokedex > ")

		input, err := a.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Fprintln(a.writer, "Internal Error:", err)
			continue
		}

		if err := a.handleInput(input); err != nil {
			fmt.Fprintln(a.writer, err)
		}
	}
}

func (a *App) handleInput(input string) error {
	parsed := repl.CleanInput(input)
	if len(parsed) == 0 {
		return nil
	}

	cmd, ok := a.commands[parsed[0]]
	if !ok {
		return fmt.Errorf("Unknown command: %s", parsed[0])
	}

	return cmd.callback(parsed[1:])
}

func (a *App) commandExit(_ []string) error {
	fmt.Fprintln(a.writer, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (a *App) commandHelp(_ []string) error {
	fmt.Fprintln(a.writer, "Welcome to the Pokedex!")
	fmt.Fprintln(a.writer, "Usage:")
	fmt.Fprintln(a.writer)

	for _, cmd := range a.commands {
		fmt.Fprintf(a.writer, "%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func (a *App) commandMap(_ []string) error {
	pageURL := a.config.next
	if pageURL == "" {
		pageURL = pokeapi.LocationAreaURL
	}

	data, err := a.client.GetLocationAreas(pageURL)
	if err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Fprintln(a.writer, area.Name)
	}

	a.config.next = data.Next
	a.config.previous = data.Previous
	return nil
}

func (a *App) commandMapb(_ []string) error {
	pageURL := a.config.previous
	if pageURL == "" {
		pageURL = pokeapi.LocationAreaURL
	}

	data, err := a.client.GetLocationAreas(pageURL)
	if err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Fprintln(a.writer, area.Name)
	}

	a.config.next = data.Next
	a.config.previous = data.Previous
	return nil
}

func (a *App) commandExplore(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <location-area-name>")
	}

	areaName := args[0]
	data, err := a.client.GetLocationArea(areaName)
	if err != nil {
		return err
	}

	fmt.Fprintf(a.writer, "Exploring %s...\n", data.Name)
	fmt.Fprintf(a.writer, "Found Pokemon:\n")
	for _, encounter := range data.PokemonEncounters {
		fmt.Fprintf(a.writer, " - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func (a *App) commandCatch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}

	pokemonName := args[0]
	data, err := a.client.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	fmt.Fprintf(a.writer, "Throwing a Pokeball at %s...\n", data.Name)

	catchThreshold := 100 - data.BaseExperience/4
	if catchThreshold < 5 {
		catchThreshold = 5
	}

	if a.catchRoll(100) >= catchThreshold {
		fmt.Fprintf(a.writer, "%s escaped!\n", data.Name)
		return nil
	}

	a.pokedex[data.Name] = data
	fmt.Fprintf(a.writer, "%s was caught!\n", data.Name)
	return nil
}

func (a *App) commandInspect(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: inspect <pokemon-name>")
	}

	pokemonName := args[0]
	data, ok := a.pokedex[pokemonName]
	if !ok {
		fmt.Fprintf(a.writer, "you have not caught that pokemon\n")
		return nil
	}

	fmt.Fprintf(a.writer, "Name: %s\n", data.Name)
	fmt.Fprintf(a.writer, "Height: %d\n", data.Height)
	fmt.Fprintf(a.writer, "Weight: %d\n", data.Weight)
	fmt.Fprintln(a.writer, "Stats:")
	for _, stat := range data.Stats {
		fmt.Fprintf(a.writer, "  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Fprintln(a.writer, "Types:")
	for _, t := range data.Types {
		fmt.Fprintf(a.writer, "  - %s\n", t.Type.Name)
	}

	return nil
}

func (a *App) commandPokedex(_ []string) error {
	fmt.Fprintln(a.writer, "Your Pokedex:")

	pokemonNames := make([]string, 0, len(a.pokedex))
	for name := range a.pokedex {
		pokemonNames = append(pokemonNames, name)
	}
	sort.Strings(pokemonNames)

	for _, name := range pokemonNames {
		fmt.Fprintf(a.writer, " - %s\n", name)
	}

	return nil
}
