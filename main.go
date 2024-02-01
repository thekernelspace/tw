package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)


func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
