package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keyboard shortcuts for the application.
type KeyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Branch   key.Binding
	Watch    key.Binding
	Quit     key.Binding
	Help     key.Binding

	Input1 key.Binding
	Input2 key.Binding
	Input3 key.Binding
	Input4 key.Binding
	Input5 key.Binding
	Input6 key.Binding
	Input7 key.Binding
	Input8 key.Binding
	Input9 key.Binding
}

// DefaultKeyMap returns the default keyboard shortcuts.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next pane")),
		ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev pane")),
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select/run")),
		Escape:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Branch:   key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "branch")),
		Watch:    key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "watch")),
		Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),

		Input1: key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "input 1")),
		Input2: key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "input 2")),
		Input3: key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "input 3")),
		Input4: key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "input 4")),
		Input5: key.NewBinding(key.WithKeys("5"), key.WithHelp("5", "input 5")),
		Input6: key.NewBinding(key.WithKeys("6"), key.WithHelp("6", "input 6")),
		Input7: key.NewBinding(key.WithKeys("7"), key.WithHelp("7", "input 7")),
		Input8: key.NewBinding(key.WithKeys("8"), key.WithHelp("8", "input 8")),
		Input9: key.NewBinding(key.WithKeys("9"), key.WithHelp("9", "input 9")),
	}
}

// ShortHelp returns a short list of key bindings for the help view.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Branch, k.Quit, k.Help}
}

// FullHelp returns the full list of key bindings for the help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.ShiftTab, k.Up, k.Down},
		{k.Enter, k.Escape, k.Branch, k.Watch},
		{k.Input1, k.Input2, k.Input3},
		{k.Quit, k.Help},
	}
}
