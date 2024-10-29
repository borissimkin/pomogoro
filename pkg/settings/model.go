package settings

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"pomogoro/pkg/router"
	"pomogoro/pkg/settings/keybinding"
)

type kindFormItem string

const (
	toggleItem kindFormItem = "toggle"
	numberItem kindFormItem = "number"
)

type formItem struct {
	id    string
	title string
	value int
	kind  kindFormItem
}

type formMap struct {
	soundNotification formItem
	pushNotification  formItem
}

func toInt(v bool) int {
	if v {
		return 1
	}

	return 0
}

func initFormMap(settings *Settings) formMap {
	return formMap{
		soundNotification: formItem{
			id:    "sound",
			title: "Sound notification",
			value: toInt(settings.Notification.Sound),
			kind:  toggleItem,
		},
		pushNotification: formItem{
			id:    "sound",
			title: "Sound notification",
			value: toInt(settings.Notification.Push),
			kind:  toggleItem,
		},
	}
}

type ListItem interface {
	View() string
}
type Model struct {
	formMap  formMap
	settings *Settings
	cursor   int
	help     help.Model
	keymap   keybinding.KeyMap
	router   *router.Router
}

func (m *Model) listItems() []formItem {
	return []formItem{
		m.formMap.soundNotification,
		m.formMap.pushNotification,
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Back):
			return m.router.To("pomodoro")
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keymap.Down):
			if m.cursor < len(m.listItems())-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func toggleView(item formItem) string {
	s := ""

	value := "off"

	if item.value == 1 {
		value = "on"
	}

	s += fmt.Sprintf("%s %s", value, item.title)

	return s
}

func (m *Model) View() string {
	s := "Settings"

	s += "\n\n"

	for index, listItem := range m.listItems() {
		cursor := " "

		if index == m.cursor {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, toggleView(listItem))
	}

	s += m.help.View(m.keymap)

	return s
}

func NewModel(r *router.Router) *Model {
	settings := NewSettings()

	return &Model{
		formMap:  initFormMap(settings),
		settings: settings,
		keymap:   keybinding.InitKeys(),
		help:     help.New(),
		router:   r,
	}
}
