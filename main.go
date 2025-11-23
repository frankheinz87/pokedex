package main

import (
	"time"

	"github.com/frankheinz87/pokedex/internal/pokeapi"
	"github.com/frankheinz87/pokedex/repl"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &repl.Config{
		PokeapiClient: pokeClient,
	}
	repl.StartRepl(cfg)
}
