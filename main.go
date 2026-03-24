package main

import (
	"fmt"
	"os"

	"emagjby/boot-pokedex/internal/pokedex"
)

func main() {
	app := pokedex.New(os.Stdin, os.Stdout)
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
