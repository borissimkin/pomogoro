package keybinding

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
)

const DefaultStepMinutes = 1

type KeyMap struct {
	Start key.Binding
	Stop  key.Binding
	Reset key.Binding
	Next  key.Binding
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Start,
		k.Stop,
		k.Reset,
		k.Next,
		k.Left,
		k.Right,
		k.Up,
		k.Down,
		k.Help,
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Start, k.Stop, k.Reset, k.Next},
		{k.Help, k.Quit},
	}
}

func InitKeys() KeyMap {
	return KeyMap{
		Start: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("␣", "start"),
		),
		Stop: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("␣", "stop"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r", "к"),
			key.WithHelp("r", "reset"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "й", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Next: key.NewBinding(
			key.WithKeys("n", "т"),
			key.WithHelp("n", "next"),
		),
		Help: key.NewBinding(
			key.WithKeys("/", "?"),
			key.WithHelp("?", "help"),
		),

		Up: key.NewBinding(
			key.WithKeys("up", "k", "л", "w", "ц"),
			key.WithHelp("↑/w/k", fmt.Sprintf("+%v min", DefaultStepMinutes)),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j", "о", "s", "ы"),
			key.WithHelp("↓/s/j", fmt.Sprintf("-%v min", DefaultStepMinutes)),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h", "р", "a", "ф"),
			key.WithHelp("←/a/h", "to left session"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l", "д", "d", "в"),
			key.WithHelp("→/d/l", "to right session"),
		),
	}
}
