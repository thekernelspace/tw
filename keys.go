package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Quit       key.Binding
	JumpTop    key.Binding
	JumpBottom key.Binding
	Expand     key.Binding
	Collapse   key.Binding
	Edit       key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc/ctrl+c", "quit"),
	),
	JumpTop: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "jump to top"),
	),
	JumpBottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "jump to bottom"),
	),
	Expand: key.NewBinding(
		key.WithKeys("enter", "right", "l"),
		key.WithHelp("enter", "expand"),
	),
	Collapse: key.NewBinding(
		key.WithKeys("backspace", "left", "h"),
		key.WithHelp("backspace", "collapse"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
}
