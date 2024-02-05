package main

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cursor           int
	currentCursorDir *Dirent
	root             Dirent
	sniper           bool
	sniperKeyBuffer  string // the currently pressed keys; "" if sniper mode is off or no keys are pressed, else the keys pressed so far
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitModel(dirpath string) Model {
	fi := getFi(dirpath)
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
		switch {
		case m.sniper:
			if msg.String() == "enter" && m.sniperKeyBuffer != "" {
				// as in go to the dirent that matches the key buffer
				dirent := getDirentFromKey(m.sniperKeyBuffer)
				if dirent != nil {
					m.currentCursorDir = dirent.Parent
					m.cursor = dirent.PosInParent
					log.Printf("sniper mode: moving to dirent: %v\n", dirent.Path())
					m.QuitSniper()
				} else {
					log.Printf("sniper mode: no dirent found for key: %v\n", m.sniperKeyBuffer)
				}
				m.QuitSniper()
				break
			}
			m.sniperKeyBuffer += msg.String()
			dirent := getDirentFromKey(m.sniperKeyBuffer)
			if dirent != nil {
				if dirent.IsDir() && dirent.Expanded {
					// stall for the next hit to tab into any dirent in the expanded dir
					break
				}
				m.currentCursorDir = dirent.Parent
				m.cursor = dirent.PosInParent
				log.Printf("sniper mode: moving to dirent: %v\n", dirent.Path())
			} else {
				log.Printf("sniper mode: no dirent found for key: %v\n", msg.String())
			}
			m.QuitSniper()
		default:
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Up):
			m.HandleMoveUp()
		case key.Matches(msg, DefaultKeyMap.Down):
			m.HandleMoveDown()
		case key.Matches(msg, DefaultKeyMap.JumpTop):
			m.HandleJumpToTop()
		case key.Matches(msg, DefaultKeyMap.JumpBottom):
			m.HandleJumpToBottom()
		case key.Matches(msg, DefaultKeyMap.Expand): // expand directory if l or right, or toggle if enter; open file if file only if enter
			if m.getCurrentDirent().IsFile() && msg.String() == "enter" {
				model, cmd := m.HandleEdit()
				if cmd != nil {
					return model, cmd
				}
			} else if m.getCurrentDirent().IsDir() {
				m.HandleExpand(msg)
			}
		case key.Matches(msg, DefaultKeyMap.Collapse): // collapse directory
			m.HandleCollapse()
		case key.Matches(msg, DefaultKeyMap.Edit): // open file with $EDITOR; ignore if it's a directory
			model, cmd := m.HandleEdit()
			if cmd != nil {
				return model, cmd
			}
			if model != nil {
				m = *model
			}
		case key.Matches(msg, DefaultKeyMap.Sniper):
			m.ToggleSniper()
		}
	}
	log.Printf("current cursor dir: %v, cursor: %v; direntpath: %v\n", m.currentCursorDir.Path(), m.cursor, m.getCurrentDirent().Path())
	return m, nil
}

func (m *Model) HandleEdit() (*Model, tea.Cmd) {
	if m.getCurrentDirent().IsFile() {
		c := EditorOpenFile(m.getCurrentDirent().Path())
		cmd := tea.ExecProcess(c, func(err error) tea.Msg {
			return nil
		})
		return m, cmd
	}
	return m, nil
}

func (m *Model) HandleExpand(msg tea.KeyMsg) {
	currentDirent := m.getCurrentDirent()
	if msg.String() == "l" && currentDirent.IsDir() && currentDirent.Expanded {
		return // l = expand only
	}
	if len(currentDirent.Dirents) == 0 {
		currentDirent.LoadDirents()
	}
	if currentDirent.IsDir() && len(currentDirent.Dirents) > 0 {
		// Expand directory
		currentDirent.Expanded = !currentDirent.Expanded
	}
}

func (m *Model) HandleCollapse() {
	currentDirent := m.getCurrentDirent()
	if currentDirent.IsDir() && currentDirent.Expanded {
		currentDirent.Expanded = false
	}
}

func (m *Model) HandleMoveUp() {
	// case 1: move out from top of subdir
	if m.cursorTopCurrentDir() {
		if m.currentCursorDir.Parent == nil {
			log.Printf("top of root; going to bottom...\n")
			m.cursor = len(m.currentCursorDir.Dirents) - 1
			return
		}
		log.Printf("moving out of %v to %v on pos %v\n", m.currentCursorDir.Path(), m.currentCursorDir.Parent.Path(), m.currentCursorDir.PosInParent)
		m.cursor = m.currentCursorDir.PosInParent
		m.currentCursorDir = m.currentCursorDir.Parent
		log.Printf("new cursor: %v\n", m.cursor)
		return
	}

	// case 2: move in to the above dirent
	direntAbove := m.currentCursorDir.Dirents[m.cursor-1]
	if direntAbove.IsDir() && direntAbove.Expanded {
		// find the bottom one in this dir
		target := direntAbove.getBottomDirent()
		m.cursor = target.PosInParent
		m.currentCursorDir = target.Parent
		log.Printf("moving in to %v; cursor new: %v\n", m.currentCursorDir, m.cursor)
		return
	}

	// in middle of the current directory
	m.cursor--
}

func (d Dirent) getClosestSuccessor() *Dirent {
	if d.PosInParent+1 >= len(d.Parent.Dirents) {
		return d.Parent.getClosestSuccessor()
	}
	return &d.Parent.Dirents[d.PosInParent+1]
}

func (m *Model) HandleMoveDown() {
	// case 1: move in; cursor moves from parent -> an expanded subdir -- when the cursor is at the subdir
	if m.getCurrentDirent().IsDir() && m.getCurrentDirent().Expanded {
		m.currentCursorDir = m.getCurrentDirent()
		m.cursor = 0
		// log.Printf("current cursor dir: %v, cursor: %v", m.currentCursorDir, m.cursor)
		return
	}

	// case 2: move out; of a subdir to the parent -- when the cursor is at the bottom of the subdir
	if m.cursorBottomCurrentDir() {
		// if there's no parent to move out to i.e you're the bottom dirent from the root, do nothing
		// TODO: cache this getBottomDirent() call
		if m.root.getBottomDirent().Equals(*m.getCurrentDirent()) {
			log.Println("bottom of root; going to top...")
			m.cursor = 0
			return
		}
		// move out to the closest ancestor with a posInParent + 1 that's still in bound
		closestSucc := m.currentCursorDir.getClosestSuccessor()
		m.cursor = closestSucc.PosInParent
		m.currentCursorDir = closestSucc.Parent
		return
	}

	// case 3: move within
	m.cursor++
}

func (m *Model) HandleJumpToBottom() {
	m.cursor = len(m.currentCursorDir.Dirents) - 1
}

func (m *Model) HandleJumpToTop() {
	m.cursor = 0
}

func (m Model) View() string {
	return m.root.Print(m)
}
