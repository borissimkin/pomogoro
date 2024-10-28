package keybinding

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Back  key.Binding
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Left,
		k.Right,
		k.Up,
		k.Down,
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back},
		{k.Up, k.Down, k.Left, k.Right},
	}
}

func InitKeys() KeyMap {
	return KeyMap{
		Back: key.NewBinding(
			key.WithKeys("esc", "b", "и"),
			key.WithHelp("esc/b", "BACK"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "й", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k", "л", "w", "ц"),
			key.WithHelp("↑/w/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j", "о", "s", "ы"),
			key.WithHelp("↓/s/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h", "р", "a", "ф"),
			key.WithHelp("←/a/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l", "д", "d", "в"),
			key.WithHelp("→/d/l", "right"),
		),
	}
}
