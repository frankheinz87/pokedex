package repl

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/frankheinz87/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

/*type Area struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Location struct {
	Name       string      `json:"name"`
	Encounters []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
}

type Response struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Area  `json:"results"`
}
*/

type Config struct {
	PokeapiClient    pokeapi.Client
	NextLocationsURL *string
	PrevLocationsURL *string
	CaughtPokemon    map[string]pokeapi.Pokemon
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    CommandExit,
	},
	"help": {
		name:        "help",
		description: "Display usage information of the Pokedex",
		callback:    CommandHelp,
	},
	"map": {
		name:        "map",
		description: "Displays the names of the next 20 location areas in the Pokemon world",
		callback:    CommandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays the names of the previous 20 location areas in the Pokemon world",
		callback:    CommandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Displays all pokemon to be encountered in the chosen area",
		callback:    CommandExp,
	},
	"catch": {
		name:        "catch",
		description: "Attempting to catch the chosen pokemon in the chosen area",
		callback:    CommandCat,
	},
	"inspect": {
		name:        "inspect",
		description: "Shows the stats of a caught pokemon",
		callback:    CommandIns,
	},
}

func CommandExit(cfg *Config, loc string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("program exited via user command")
}

func CommandHelp(cfg *Config, loc string) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
	
help: Displays a help message
exit: Exit the Pokedex`)
	return nil
}

func CommandExp(cfg *Config, loc string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.NextLocationsURL != nil {
		url = *cfg.NextLocationsURL
	}

	locResp, err := cfg.PokeapiClient.LocationPokemon(url, loc)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", loc)
	fmt.Println("Found Pokemon:")

	for _, encounter := range locResp.Encounters {
		fmt.Println(encounter.Pokemon.Name)

	}

	return nil
}

func CommandCat(cfg *Config, poc string) error {
	url := "https://pokeapi.co/api/v2/pokemon"
	if cfg.NextLocationsURL != nil {
		url = *cfg.NextLocationsURL
	}

	locResp, err := cfg.PokeapiClient.GetPokemon(url, poc)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", poc)
	roll := rand.Intn(locResp.Base_experience)
	if roll > 40 {
		fmt.Printf("%s excaped!\n", poc)
	} else {
		cfg.CaughtPokemon[poc] = locResp
		fmt.Printf("%s was caught\n", poc)
	}
	return nil
}

func CommandMap(cfg *Config, loc string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.NextLocationsURL != nil {
		url = *cfg.NextLocationsURL
	}

	locResp, err := cfg.PokeapiClient.ListLocations(url)
	if err != nil {
		return err
	}

	cfg.NextLocationsURL = locResp.Next
	cfg.PrevLocationsURL = locResp.Previous

	for _, area := range locResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func CommandMapb(cfg *Config, loc string) error {
	if cfg.PrevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url := *cfg.PrevLocationsURL

		locResp, err := cfg.PokeapiClient.ListLocations(url)
		if err != nil {
			return err
		}

		cfg.NextLocationsURL = locResp.Next
		cfg.PrevLocationsURL = locResp.Previous

		for _, area := range locResp.Results {
			fmt.Println(area.Name)
		}
		return nil
	}
}

func CommandIns(cfg *Config, poc string) error {
	_, ok := cfg.CaughtPokemon[poc]
	if !ok {
		fmt.Printf("%s has not been caught yet!\n", poc)
		return nil
	}
	fmt.Printf("Name: %s\n", poc)
	fmt.Printf("Height: %v\n", cfg.CaughtPokemon[poc].Height)
	fmt.Printf("Weight: %v\n", cfg.CaughtPokemon[poc].Weight)
	fmt.Println("Stats:")
	for _, stat := range cfg.CaughtPokemon[poc].Stats {
		fmt.Printf("	-%s: %v\n", stat.Stat.Name, stat.Base_stat)
	}
	fmt.Println("Types:")
	for _, t := range cfg.CaughtPokemon[poc].Types {
		fmt.Printf("	-%s\n", t.Type.Name)
	}
	return nil
}

/*func Execute(cfg *Config, words []string) error {
	if len(words) == 0 {
		return nil
	}

	cmdName := words[0]
	cmd, exists := commands[cmdName]
	if !exists {
		return fmt.Errorf("Unknown command")
	}

	return cmd.callback(cfg)
}

func decodeJSONResponse(res *http.Response) ([]Area, error) {
	var areas []Area
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&areas); err != nil {
		return nil, err
	}
	return areas, nil
}*/

func CleanInput(text string) []string {
	return strings.Fields(strings.TrimSpace(strings.ToLower(text)))

}

func StartRepl(cfg *Config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !reader.Scan() {
			return
		}

		words := CleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		Location := ""
		if len(words) > 1 {
			Location = words[1]
		}

		cmd, exists := getCommands()[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(cfg, Location); err != nil {
			fmt.Println(err)
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    CommandExit,
		},
		"help": {
			name:        "help",
			description: "Display usage information of the Pokedex",
			callback:    CommandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world",
			callback:    CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world",
			callback:    CommandMapb,
		},
		"explore": {
			name:        "mapb",
			description: "Displays all pokemon to be encountered in the chosen area",
			callback:    CommandExp,
		},
		"catch": {
			name:        "catch",
			description: "Attempting to catch the chosen pokemon in the chosen area",
			callback:    CommandCat,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows the stats of a caught pokemon",
			callback:    CommandIns,
		},
	}
}
