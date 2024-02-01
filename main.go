package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)


func main() {
	p := tea.NewProgram(InitModel())
	if _, err := p.Run(); err != nil {
		log.Panicf("Error: %v\n", err)
		os.Exit(1)
	}
}
