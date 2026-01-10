package main

import (
	"math/rand"
	"time"

	"github.com/frankheinz87/pokedex/internal/pokeapi"
	"github.com/frankheinz87/pokedex/repl"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	catchRegister := make(map[string]pokeapi.Pokemon)
	cfg := &repl.Config{
		PokeapiClient: pokeClient,
		CaughtPokemon: catchRegister,
	}
	repl.StartRepl(cfg)
}
