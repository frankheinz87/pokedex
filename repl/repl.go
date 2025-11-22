package repl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Area struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Response struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Area  `json:"results"`
}

type Config struct {
	NextLocationsURL *string
	PrevLocationsURL *string
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
}

func CommandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("program exited via user command")
}

func CommandHelp(cfg *Config) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
	
help: Displays a help message
exit: Exit the Pokedex`)
	return nil
}

func CommandMap(cfg *Config) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.NextLocationsURL != nil {
		url = *cfg.NextLocationsURL
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var locResp Response
	err = json.Unmarshal(data, &locResp)
	if err != nil {
		log.Fatal(err)
	}

	cfg.NextLocationsURL = locResp.Next
	cfg.PrevLocationsURL = locResp.Previous

	for _, area := range locResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func CommandMapb(cfg *Config) error {
	if cfg.PrevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url := *cfg.PrevLocationsURL

		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var locResp Response
		err = json.Unmarshal(data, &locResp)
		if err != nil {
			log.Fatal(err)
		}

		cfg.NextLocationsURL = locResp.Next
		cfg.PrevLocationsURL = locResp.Previous

		for _, area := range locResp.Results {
			fmt.Println(area.Name)
		}
		return nil
	}
}

func Execute(cfg *Config, words []string) error {
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
}

func CleanInput(text string) []string {
	return strings.Fields(strings.TrimSpace(strings.ToLower(text)))

}
