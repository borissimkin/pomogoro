package keybinding

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Reset key.Binding
	Enter key.Binding
	Back  key.Binding
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Enter,
		k.Back,
		k.Left,
		k.Right,
		k.Up,
		k.Down,
		k.Reset,
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
		Reset: key.NewBinding(
			key.WithKeys("r", "к"),
			key.WithHelp("r", "reset all to defaults")),
		Enter: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "toggle"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "b", "и"),
			key.WithHelp("esc/b", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "й", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k", "л", "w", "ц"),
			key.WithHelp("↑/w/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j", "о", "s", "ы"),
			key.WithHelp("↓/s/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h", "р", "a", "ф"),
			key.WithHelp("←/a/h", "decrease"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l", "д", "d", "в"),
			key.WithHelp("→/d/l", "increase"),
		),
	}
}
