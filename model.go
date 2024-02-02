package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cursor           int
	currentCursorDir *Dirent
	root             Dirent
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitModel() Model {
	// get the directory entries for the current directory
	dir, err := os.Open(".")
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}

	defer dir.Close()

	fi, err := dir.Stat()
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}

	root := Dirent{
		fi:       fi,
		Level:    0,
		Expanded: true,
	}
	root.LoadDirents()

	return Model{
		currentCursorDir: &root,
		cursor:           0,
		root:             root,
	}
}

func (m Model) getCurrentDirent() *Dirent {
	return &m.currentCursorDir.Dirents[m.cursor]
}

func (m Model) cursorBottomCurrentDir() bool {
	return m.cursor == len(m.currentCursorDir.Dirents)-1
}

func (m Model) cursorTopCurrentDir() bool {
	return m.cursor == 0
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			// case 1: move out from top of subdir
			if m.cursorTopCurrentDir() {
				if m.currentCursorDir.Parent == nil {
					log.Printf("top of root; doing nothing...\n")
					break
				}
				log.Printf("moving out of %v to %v\n", m.currentCursorDir.Path(), m.currentCursorDir.Parent.Path())
				m.cursor = m.currentCursorDir.PosInParent
				m.currentCursorDir = m.currentCursorDir.Parent
				log.Printf("new cursor: %v\n", m.cursor)
				break
			}

			// case 2: move within from the next dirent it's pointing to
			direntAboveCurrent := m.currentCursorDir.Dirents[m.cursor-1]
			if direntAboveCurrent.IsDir() && direntAboveCurrent.Expanded {
				log.Printf("moving in to %v; cursor new: %v\n", direntAboveCurrent.Path(), len(direntAboveCurrent.Dirents)-1)
				m.cursor = len(direntAboveCurrent.Dirents) - 1
				m.currentCursorDir = &direntAboveCurrent
				break
			}

			// in middle of the current directory
			m.cursor--
		case "down", "j":
			// case 1: move in; cursor moves from parent -> an expanded subdir -- when the cursor is at the subdir
			if m.getCurrentDirent().IsDir() && m.getCurrentDirent().Expanded {
				m.currentCursorDir = m.getCurrentDirent()
				m.cursor = 0
				// log.Printf("current cursor dir: %v, cursor: %v", m.currentCursorDir, m.cursor)
				break
			}

			// case 2: move out; of a subdir to the parent -- when the cursor is at the bottom of the subdir
			if m.cursorBottomCurrentDir() {
				// if there's no parent to move out to i.e you're the bottom dirent from the root, do nothing
				// TODO: cache this getBottomDirent() call
				if m.root.getBottomDirent().Equals(*m.getCurrentDirent()) {
					log.Println("bottom of root; doing nothing...")
					break
				}
				// there's a parent so move out
				m.cursor = m.currentCursorDir.PosInParent + 1
				m.currentCursorDir = m.currentCursorDir.Parent
				break
			}

			// case 3: move within
			m.cursor++
		case "enter", "right", "l": // expand directory if l or right, or toggle if enter; open file if file
			currentDirent := m.getCurrentDirent()
			if currentDirent.IsDir() {
				// Expand directory
				currentDirent.Expanded = !currentDirent.Expanded
				// load the directory entries for the current directory
				if currentDirent.Expanded && len(currentDirent.Dirents) == 0 {
					currentDirent.LoadDirents()
				}
			} else if currentDirent.IsFile() {
				// open file with $EDITOR
				c := EditorOpenFile(currentDirent.Path())
				cmd := tea.ExecProcess(c, func(err error) tea.Msg {
					return nil
				})
				return m, cmd
			}
		case "left", "h": // collapse directory
			currentDirent := m.getCurrentDirent()
			if currentDirent.IsDir() && currentDirent.Expanded {
				currentDirent.Expanded = false
			}
		case "e": // open file with $EDITOR; ignore if it's a directory
			if m.getCurrentDirent().IsFile() {
				c := EditorOpenFile(m.getCurrentDirent().Path())
				cmd := tea.ExecProcess(c, func(err error) tea.Msg {
					return nil
				})
				return m, cmd
			}
		}
	}
	log.Printf("current cursor dir: %v, cursor: %v; direntpath: %v\n", m.currentCursorDir.Path(), m.cursor, m.getCurrentDirent().Path())
	return m, nil
}

func (m Model) View() string {
	return m.root.Print(m)
}
