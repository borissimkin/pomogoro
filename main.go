package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionType string

const (
	workSession      SessionType = "workSession"
	breakSession     SessionType = "breakSession"
	longBreakSession SessionType = "longBreakSession"
)

var workSessionSettings = SessionSettings{
	title: "Pomodoro",
}

var breakSessionSettings = SessionSettings{
	title: "Short Break",
}

var longBreakSessionSettings = SessionSettings{
	title: "Long Break",
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
			workSession:      time.Minute * 25,
			breakSession:     time.Minute * 5,
			longBreakSession: time.Minute * 15,
		},
	}
}

type SessionSettings struct {
	title string
	emoji int
	color string
}

type Pomodoro struct {
	currentSessionType SessionType
	settings           *Settings
	lastSessionType    SessionType
	completed          map[SessionType]int
}

func newPomodoro(settings *Settings) *Pomodoro {
	return &Pomodoro{
		currentSessionType: workSession,
		completed:          make(map[SessionType]int),
		settings:           settings,
	}
}

func (p *Pomodoro) getNextSessionType() SessionType {
	if p.lastSessionType == "" || p.lastSessionType == longBreakSession || p.lastSessionType == breakSession {
		return workSession
	}

	if p.completed[workSession]%p.settings.workSessionsUntilLongBreak != 0 {
		return workSession
	}

	return breakSession
}

//func (p *Pomodoro)  {}

const timeout = time.Minute * 25

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
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = timeout
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

func (m model) View() string {
	// For a more detailed timer view you could read m.timer.Timeout to get
	// the remaining time as a time.Duration and skip calling m.timer.View()
	// entirely.
	s := m.timer.View()

	if m.timer.Timedout() {
		s = "All done!"
	}
	s += "\n"
	if !m.quitting {
		s = "Exiting in " + s
		s += m.helpView()
	}
	return s
}

func main() {
	settings := newSettings()

	pomodoro := newPomodoro(settings)

	m := model{
		timer:    timer.NewWithInterval(timeout, time.Second),
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
