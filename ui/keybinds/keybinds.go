package keybinds

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Place key.Binding
	Enter key.Binding
	Help  key.Binding
	Quit  key.Binding
	Esc   key.Binding
	// For help display
	// Move key.Binding
}

// // ShortHelp returns keybindings to be shown in the mini help view. It's part
// // of the key.Map interface.
// func (k keyMap) ShortHelp() []key.Binding {
// 	return []key.Binding{k.Move, k.Place, k.Pause, k.Help, k.Quit}
// }

// // FullHelp returns keybindings for the expanded help view. It's part of the
// // key.Map interface.
// func (k keyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{k.Up, k.Down, k.Left, k.Right}, // first column
// 		{k.Help, k.Quit},                // second column
// 	}
// }

var KeyBinds = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h", "a"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l", "d"),
		key.WithHelp("→/l", "move right"),
	),
	Place: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("<space>", "place"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "pause"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("<esc>", "esc"),
	),
}
