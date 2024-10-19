package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"pomogoro/packages/keybinding"
	"sort"
	"time"
)

type model struct {
	timer    timer.Model
	keymap   keybinding.KeyMap
	help     help.Model
	quitting bool
	pomodoro *Pomodoro
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.Stop.SetEnabled(m.timer.Running())
		m.keymap.Start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		nextSession := m.pomodoro.nextSession()
		notification := m.pomodoro.sessionSettings[nextSession].notification
		notify(notification.title, notification.message)
		m.timer.Timeout = m.pomodoro.getDuration()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Reset):
			m.timer.Timeout = m.pomodoro.getDuration()
		case key.Matches(msg, m.keymap.Start, m.keymap.Stop):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.keymap.Next):
			m.pomodoro.nextSession()
			m.timer.Timeout = m.pomodoro.getDuration()
		case key.Matches(msg, m.keymap.Right):
			m.pomodoro.setSession(getSessionType(m.pomodoro.currentSessionType + 1))
			m.timer.Timeout = m.pomodoro.getDuration()
		case key.Matches(msg, m.keymap.Left):
			m.pomodoro.setSession(getSessionType(m.pomodoro.currentSessionType - 1))
			m.timer.Timeout = m.pomodoro.getDuration()
		}
	}

	return m, nil
}

func renderTotalSessions(p *Pomodoro) string {
	return fmt.Sprintf("Total Work sessions: %v", p.totalWorkSessions())
}

func renderSessionsBeforeLongBreak(p *Pomodoro) string {
	return fmt.Sprintf("Sessions left before the long break: %v", p.sessionsBeforeLongBreak())
}

func renderSessionTypes(p *Pomodoro) string {
	s := ""

	settings := make([]*SessionSettings, 0, len(p.sessionSettings))

	for _, value := range p.sessionSettings {
		settings = append(settings, value)
	}

	sort.Slice(settings, func(i, j int) bool {
		return settings[i].sessionType < settings[j].sessionType
	})

	for _, item := range settings {
		cursor := " "

		var style = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			MarginRight(1).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color(item.backgroundColor)).
			Padding(0, 1)

		if item.sessionType != p.currentSessionType {
			style = style.
				Faint(true)
		} else {
			cursor = "*"
		}

		s += style.Render(fmt.Sprintf("%s %s", cursor, item.title))
	}

	return s
}

func renderBreakLine() string {
	return "\n"
}

func renderTime(m model) string {
	var style = lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Bold(true)

	return style.Render(m.timer.View())
}

func (m model) View() string {
	s := renderSessionTypes(m.pomodoro)

	s += renderBreakLine()
	s += renderBreakLine()

	s += renderTime(m)

	s += renderBreakLine()
	s += renderBreakLine()

	s += renderTotalSessions(m.pomodoro)
	s += renderBreakLine()
	s += renderSessionsBeforeLongBreak(m.pomodoro)

	s += renderBreakLine()

	s += m.help.View(m.keymap)

	return s
}

func main() {
	settings := newSettings()

	pomodoro := newPomodoro(settings)

	m := model{
		timer:    timer.NewWithInterval(pomodoro.getDuration(), time.Second),
		pomodoro: pomodoro,
		keymap:   keybinding.InitKeys(),
		help:     help.New(),
	}
	m.keymap.Start.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
