package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const DEBUG = true

func main() {
  if DEBUG {
    f, err := os.Create("debug.log")
    if err != nil {
      log.Fatal(err)
    }
    defer f.Close()
    log.SetOutput(f)
  }

	p := tea.NewProgram(InitModel())
	if _, err := p.Run(); err != nil {
		log.Panicf("Error: %v\n", err)
		os.Exit(1)
	}
}
