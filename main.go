package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionType int

const (
	workSession      SessionType = 1
	breakSession     SessionType = 2
	longBreakSession SessionType = 3
)

var workSessionSettings = SessionSettings{
	sessionType: workSession,
	title:       "Pomodoro",
}

var breakSessionSettings = SessionSettings{
	sessionType: breakSession,
	title:       "Short Break",
}

var longBreakSessionSettings = SessionSettings{
	sessionType: longBreakSession,
	title:       "Long Break",
}

type durations map[SessionType]time.Duration

type Settings struct {
	workSessionsUntilLongBreak int
	durations                  durations
}

func (s *Settings) getDuration(sessionType SessionType) time.Duration {
	return s.durations[sessionType]
}

func newSettings() *Settings {
	return &Settings{
		workSessionsUntilLongBreak: 4,
		durations: durations{
			//workSession:      time.Minute * 25,
			//breakSession:     time.Minute * 5,
			//longBreakSession: time.Minute * 15,
			workSession:      time.Second * 10,
			breakSession:     time.Second * 2,
			longBreakSession: time.Second * 7,
		},
	}
}

type SessionSettings struct {
	sessionType SessionType
	title       string
	emoji       int
	color       string
}

type Pomodoro struct {
	currentSessionType  SessionType
	settings            *Settings
	sessionSettings     map[SessionType]*SessionSettings
	previousSessionType SessionType
	completed           map[SessionType]int
}

func newPomodoro(settings *Settings) *Pomodoro {
	return &Pomodoro{
		currentSessionType: workSession,
		completed:          make(map[SessionType]int),
		settings:           settings,
		sessionSettings: map[SessionType]*SessionSettings{
			workSession:      &workSessionSettings,
			breakSession:     &breakSessionSettings,
			longBreakSession: &longBreakSessionSettings,
		},
	}
}

func (p *Pomodoro) getDuration() time.Duration {
	return p.settings.getDuration(p.currentSessionType)
}

func (p *Pomodoro) getNextSessionType() SessionType {
	if p.previousSessionType != workSession {
		return workSession
	}

	if p.completed[workSession]%p.settings.workSessionsUntilLongBreak == 0 {
		return longBreakSession
	}

	return breakSession
}

func (p *Pomodoro) nextSession() {
	p.completed[p.currentSessionType]++
	p.previousSessionType = p.currentSessionType

	p.currentSessionType = p.getNextSessionType()
}

type model struct {
	timer    timer.Model
	keymap   keymap
	help     help.Model
	quitting bool
	pomodoro *Pomodoro
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		m.pomodoro.nextSession()
		m.timer.Timeout = m.pomodoro.getDuration()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = m.pomodoro.getDuration()
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func renderSessionTypes(m model) string {
	s := ""

	settings := make([]*SessionSettings, 0, len(m.pomodoro.sessionSettings))

	for _, value := range m.pomodoro.sessionSettings {
		settings = append(settings, value)
	}

	sort.Slice(settings, func(i, j int) bool {
		return settings[i].sessionType < settings[j].sessionType
	})

	for _, item := range settings {
		cursor := " "

		if item.sessionType == m.pomodoro.currentSessionType {
			cursor = "*"
		}

		s += fmt.Sprintf("%s %s\n", cursor, item.title)
	}

	return s
}

func renderBreakLine() string {
	return "\n"
}

func (m model) View() string {
	// For a more detailed timer view you could read m.timer.Timeout to get
	// the remaining time as a time.Duration and skip calling m.timer.View()
	// entirely.
	s := renderSessionTypes(m)

	s += renderBreakLine()

	s += m.timer.View()

	s += renderBreakLine()

	s += m.helpView()

	return s
}

func main() {
	settings := newSettings()

	pomodoro := newPomodoro(settings)

	m := model{
		timer:    timer.NewWithInterval(pomodoro.getDuration(), time.Second),
		pomodoro: pomodoro,
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
	m.keymap.start.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
