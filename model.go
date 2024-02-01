package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cursor int
	root   Dirent
}

func (m Model) Init() tea.Cmd {
	return nil
}

func initModel() Model {
	// get the directory entries for the current directory
	dir, err := os.Open(".")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	defer dir.Close()

	fi, err := dir.Stat()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	root := Dirent{
		Path:  fi.Name(),
		IsDir: fi.IsDir(),
	}
	root.LoadDirents()

	return Model{
		cursor: 0,
		root:   root,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.root.Dirents)-1 {
				m.cursor++
			}
		case "enter":
			if m.root.Dirents[m.cursor].IsDir {
				// Expand directory
			} else {
				// Open file with the default $EDITOR
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vi"
				}
				c := exec.Command(editor, m.root.Dirents[m.cursor].Path)
				cmd := tea.ExecProcess(c, func(err error) tea.Msg {
					return nil
				})
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := ""
	for i, dirent := range m.root.Dirents {
		cursor := " "
		if i == m.cursor {
			// Highlight the cursor
			cursor = ">"
		}
		// Render the dirent
		endingSlash := ""
		if dirent.IsDir {
			endingSlash = "/"
		}
		s += fmt.Sprintf("%s%s%s\n", cursor, dirent.Path, endingSlash)
	}
	return s
}

