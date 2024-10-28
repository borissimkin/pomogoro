package settings

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"pomogoro/pkg/router"
	"pomogoro/pkg/settings/keybinding"
)

type Model struct {
	help   help.Model
	keymap keybinding.KeyMap
	router *router.Router
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Back):
			return m.router.To("pomodoro")
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	s := "SETTINGS VIEW! =)"

	s += m.help.View(m.keymap)

	return s
}

func NewModel(r *router.Router) *Model {
	return &Model{
		keymap: keybinding.InitKeys(),
		help:   help.New(),
		router: r,
	}
}
