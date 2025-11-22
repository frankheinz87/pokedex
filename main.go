package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/frankheinz87/pokedex/repl"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &repl.Config{
		NextLocationsURL: nil,
		PrevLocationsURL: nil,
	}
	for i := 1; i > 0; i++ {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			words := repl.CleanInput(text)
			//fmt.Println("Your command was:", words[0])
			if err := repl.Execute(cfg, words); err != nil {
				fmt.Println(err)
			}
		}
	}

}
