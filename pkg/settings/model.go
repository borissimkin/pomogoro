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
	"time"
)

type kindFormItem string

const (
	toggleItem kindFormItem = "toggle"
	numberItem kindFormItem = "number"
)

const (
	minLimit = 1
	maxLimit = 9999
)

type limits struct {
	min int
	max int
}

type formItem struct {
	title  string
	value  int
	kind   kindFormItem
	limits *limits
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

	value := item.value + 1

	if item.limits != nil && item.limits.max < value {
		return
	}

	item.value = value
}

func (item *formItem) decrease() {
	if item.isToggle() {
		item.value = 0
	}

	value := item.value - 1

	if item.limits != nil && item.limits.min > value {
		return
	}

	item.value = value
}

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

func toBool(v int) bool {
	return v == 1
}

func initFormMap(settings *Settings) formMap {
	return formMap{
		workMinutes: &formItem{
			title: "minutes: Pomodoro",
			value: int(settings.Durations[session.Work].Minutes()),
			kind:  numberItem,
			limits: &limits{
				min: minLimit,
				max: maxLimit,
			},
		},
		breakMinutes: &formItem{
			title: "minutes: Break",
			value: int(settings.Durations[session.Break].Minutes()),
			kind:  numberItem,
			limits: &limits{
				min: minLimit,
				max: maxLimit,
			},
		},
		longBreakMinutes: &formItem{
			title: "minutes: Long Break",
			value: int(settings.Durations[session.LongBreak].Minutes()),
			kind:  numberItem,
			limits: &limits{
				min: minLimit,
				max: maxLimit,
			},
		},
		workSessionsBeforeLongBreak: &formItem{
			title: "Long Break interval",
			value: settings.WorkSessionsUntilLongBreak,
			kind:  numberItem,
			limits: &limits{
				min: 0,
				max: maxLimit,
			},
		},
		workAutoStart: &formItem{
			title: "Auto start: Pomodoro",
			value: toInt(settings.AutoStart[session.Work]),
			kind:  toggleItem,
		},
		breakAutoStart: &formItem{
			title: "Auto start: Break",
			value: toInt(settings.AutoStart[session.Break]),
			kind:  toggleItem,
		},
		longBreakAutoStart: &formItem{
			title: "Auto start: Long Break",
			value: toInt(settings.AutoStart[session.LongBreak]),
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

type Model struct {
	formMap  formMap
	settings *Settings
	cursor   int
	help     help.Model
	keymap   keybinding.KeyMap
	router   *router.Router
}

func (m *Model) resetSettings() {
	settings := DefaultSettings()

	m.settings = &settings
	m.formMap = initFormMap(&settings)
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
		case key.Matches(msg, m.keymap.Reset):
			m.resetSettings()
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keymap.Back):
			m.save()
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

func mapToSettings(form formMap) Settings {
	return Settings{
		WorkSessionsUntilLongBreak: form.workSessionsBeforeLongBreak.value,
		Durations: durations{
			session.Work:      time.Minute * time.Duration(form.workMinutes.value),
			session.Break:     time.Minute * time.Duration(form.breakMinutes.value),
			session.LongBreak: time.Minute * time.Duration(form.longBreakMinutes.value),
		},
		ShowProgressBar: toBool(form.showProgressBar.value),
		Notification: Notification{
			Sound: toBool(form.soundNotification.value),
			Push:  toBool(form.pushNotification.value),
		},
		AutoStart: AutoStart{
			session.Work:      toBool(form.workAutoStart.value),
			session.Break:     toBool(form.breakAutoStart.value),
			session.LongBreak: toBool(form.longBreakAutoStart.value),
		},
	}
}

func (m *Model) save() {
	settings := mapToSettings(m.formMap)

	_ = newStorage().Save(settings)
}

func (m *Model) Init() tea.Cmd {
	return tea.ClearScreen
}

var (
	onStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	offStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
)

func toggleItemView(item *formItem) string {
	s := ""

	value := offStyle.Render("off")

	if item.value == 1 {
		value = onStyle.Render("on")
	}

	s += fmt.Sprintf("%s %s", value, item.title)

	return s
}

func numberItemView(item *formItem) string {
	if item.value <= 0 {
		return fmt.Sprintf("%s %s", offStyle.Render("None"), item.title)
	}

	return fmt.Sprintf("%v %s", item.value, item.title)
}

func (m *Model) View() string {
	s := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		MarginLeft(2).
		PaddingLeft(1).
		PaddingRight(1).
		Render("Settings")

	s += "\n"

	for index, listItem := range m.listItems() {
		cursor := " "
		view := ""

		if index == m.cursor {
			cursor = ">"
		}

		if listItem.isToggle() {
			view = toggleItemView(listItem)
		} else if listItem.isNumber() {
			view = numberItemView(listItem)
		}

		s += fmt.Sprintf("%s %s\n", cursor, view)
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
