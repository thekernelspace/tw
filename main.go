package main

import (
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var DEBUG = os.Getenv("DEBUG") != ""

func main() {
	if DEBUG {
		f, err := os.Create("debug.log")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println("Debug mode enabled")
	} else {
		log.SetOutput(io.Discard)
	}

	// load the config
	LoadGlobalCfg()

	p := tea.NewProgram(InitModel())
	if _, err := p.Run(); err != nil {
		log.Panicf("Error: %v\n", err)
		os.Exit(1)
	}
}
