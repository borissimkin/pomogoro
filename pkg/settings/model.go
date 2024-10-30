package settings

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pomogoro/pkg/router"
	"pomogoro/pkg/session"
	"pomogoro/pkg/settings/keybinding"
)

// todo: reset settings

type kindFormItem string

const (
	toggleItem kindFormItem = "toggle"
	numberItem kindFormItem = "number"
)

type formItem struct {
	title string
	value int
	kind  kindFormItem
}

func (item *formItem) isToggle() bool {
	return item.kind == toggleItem
}

func (item *formItem) isNumber() bool {
	return item.kind == numberItem
}

func (item *formItem) enter() {
	if !item.isToggle() {
		return
	}

	if item.value >= 1 {
		item.value = 0
	} else {
		item.value = 1
	}
}

func (item *formItem) increase() {
	if item.isToggle() {
		item.value = 1
		return
	}

	item.value++
}

func (item *formItem) decrease() {
	if item.isToggle() {
		item.value = 0
	}

	item.value--
}

//- [ ] устанавливать количества минут для всех типов сессий
//- [ ] тогл автовспроизведения следующей сессии
//- [ ] количество рабочих сессий до длинного отдыха
//- [ ] тогл оповещения
//- [ ] отдельный тогл на звук
//- [ ] сохранять настройки в файл и подгружать при старте
//- [ ] показывать прогресс бар?

type formMap struct {
	workMinutes                 *formItem
	breakMinutes                *formItem
	longBreakMinutes            *formItem
	workSessionsBeforeLongBreak *formItem
	workAutoStart               *formItem
	breakAutoStart              *formItem
	longBreakAutoStart          *formItem
	soundNotification           *formItem
	pushNotification            *formItem
	showProgressBar             *formItem
}

func toInt(v bool) int {
	if v {
		return 1
	}

	return 0
}

func initFormMap(settings *Settings) formMap {
	test := &formItem{
		title: "minutes: Pomodoro",
		value: int(settings.Durations[session.Work].Minutes()),
		kind:  numberItem,
	}

	return formMap{
		workMinutes: test,
		breakMinutes: &formItem{
			title: "minutes: Break",
			value: int(settings.Durations[session.Break].Minutes()),
			kind:  numberItem,
		},
		longBreakMinutes: &formItem{
			title: "minutes: Long Break",
			value: int(settings.Durations[session.LongBreak].Minutes()),
			kind:  numberItem,
		},
		workSessionsBeforeLongBreak: &formItem{
			title: "Long Break interval",
			value: settings.WorkSessionsUntilLongBreak,
			kind:  numberItem,
		},
		workAutoStart: &formItem{
			title: "Auto start: Pomodoro",
			value: toInt(settings.AutoStart.Work),
			kind:  toggleItem,
		},
		breakAutoStart: &formItem{
			title: "Auto start: Break",
			value: toInt(settings.AutoStart.Break),
			kind:  toggleItem,
		},
		longBreakAutoStart: &formItem{
			title: "Auto start: Long Break",
			value: toInt(settings.AutoStart.LongBreak),
			kind:  toggleItem,
		},
		soundNotification: &formItem{
			title: "Sound notification",
			value: toInt(settings.Notification.Sound),
			kind:  toggleItem,
		},
		pushNotification: &formItem{
			title: "Push notification",
			value: toInt(settings.Notification.Push),
			kind:  toggleItem,
		},
		showProgressBar: &formItem{
			title: "Show progress bar",
			value: toInt(settings.ShowProgressBar),
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

func (m *Model) listItems() []*formItem {
	return []*formItem{
		m.formMap.workMinutes,
		m.formMap.breakMinutes,
		m.formMap.longBreakMinutes,
		m.formMap.workSessionsBeforeLongBreak,
		m.formMap.workAutoStart,
		m.formMap.breakAutoStart,
		m.formMap.longBreakAutoStart,
		m.formMap.soundNotification,
		m.formMap.pushNotification,
		m.formMap.showProgressBar,
	}
}

func (m *Model) currentItem() *formItem {
	return m.listItems()[m.cursor]
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Back):
			return m.router.To("pomodoro")
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Enter):
			m.currentItem().enter()
		case key.Matches(msg, m.keymap.Left):
			m.currentItem().decrease()
		case key.Matches(msg, m.keymap.Right):
			m.currentItem().increase()
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

var (
	onStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	offStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
)

func toggleView(item *formItem) string {
	s := ""

	value := offStyle.Render("off")

	if item.value == 1 {
		value = onStyle.Render("on")
	}

	s += fmt.Sprintf("%s %s", value, item.title)

	return s
}

func (m *Model) View() string {
	s := "  Settings"

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
