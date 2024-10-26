package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"pomogoro/packages/keybinding"
	"pomogoro/packages/sound"
	"sort"
	"strings"
	"time"
)

type model struct {
	progress    progress.Model
	timer       timer.Model
	initTime    time.Duration
	soundPlayer *sound.Player
	keymap      keybinding.KeyMap
	help        help.Model
	pomodoro    *Pomodoro
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func getIncreasedTime(timeout time.Duration, minutes time.Duration) time.Duration {
	return timeout + minutes*time.Minute
}

func getDecreasedTime(timeout time.Duration, minutes time.Duration) time.Duration {
	return timeout - minutes*time.Minute
}

func setTime(m *model, duration time.Duration) {
	m.timer.Timeout = duration
	m.initTime = duration
}

const (
	progressBarMaxWidth = 43
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

		m.progress.Width = msg.Width
		if m.progress.Width > progressBarMaxWidth {
			m.progress.Width = progressBarMaxWidth
		}
		return m, nil

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
		m.soundPlayer.Play()
		setTime(&m, m.pomodoro.getDuration())
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Reset):
			setTime(&m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Start, m.keymap.Stop):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.keymap.Next):
			m.pomodoro.nextSession()
			setTime(&m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Right):
			m.pomodoro.setSession(getSessionType(m.pomodoro.currentSessionType + 1))
			setTime(&m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Left):
			m.pomodoro.setSession(getSessionType(m.pomodoro.currentSessionType - 1))
			setTime(&m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Up):
			newTimeout := getIncreasedTime(m.timer.Timeout, keybinding.DefaultStepMinutes)
			m.initTime += keybinding.DefaultStepMinutes * time.Minute
			m.timer.Timeout = newTimeout
		case key.Matches(msg, m.keymap.Down):
			newTimeout := getDecreasedTime(m.timer.Timeout, keybinding.DefaultStepMinutes)
			m.initTime -= keybinding.DefaultStepMinutes * time.Minute
			if newTimeout < 0 {
				m.pomodoro.nextSession()
				setTime(&m, m.pomodoro.getDuration())
			} else {
				m.timer.Timeout = newTimeout
			}
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

// todo: pass count
func renderBreakLine() string {
	return "\n"
}

func removeMilliseconds(time string) string {
	zero := "0s"

	if time == zero {
		return time
	}

	pointer := "."

	part := strings.Split(time, pointer)[0]

	hasMs := strings.Contains(part, "ms")

	if strings.Contains(time, pointer) {
		return part + "s"
	}

	if hasMs {
		return zero
	}

	return time
}

func renderTime(m model) string {
	var style = lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Bold(true)

	if !m.timer.Running() {
		style = style.Faint(true)
	}

	// remove milliseconds cause interval with time.Second has bug
	return style.Render(removeMilliseconds(m.timer.View()))
}

func getPercent(m *model) float64 {
	return float64(m.initTime-m.timer.Timeout) / float64(m.initTime)
}

func isPause(m *model) bool {
	return m.keymap.Start.Enabled()
}

func renderProgressBar(m *model) string {
	color := m.pomodoro.currentSessionSettings().backgroundColor

	if isPause(m) {
		color = "#4b4453" // todo
	}

	m.progress.FullColor = color

	return m.progress.ViewAs(getPercent(m))
}

func (m model) View() string {
	s := renderSessionTypes(m.pomodoro)

	s += renderBreakLine()
	s += renderBreakLine()

	s += renderTime(m)

	s += renderBreakLine()

	s += renderProgressBar(&m)

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
	soundPlayer := sound.NewSoundPlayer()

	soundPlayer.InitSoundContext()
	pomodoro := newPomodoro(settings)

	initTime := pomodoro.getDuration()

	m := model{
		progress:    progress.New(progress.WithSolidFill(pomodoro.currentSessionSettings().backgroundColor), progress.WithoutPercentage()),
		timer:       timer.NewWithInterval(initTime, time.Millisecond),
		initTime:    initTime,
		pomodoro:    pomodoro,
		soundPlayer: soundPlayer,
		keymap:      keybinding.InitKeys(),
		help:        help.New(),
	}
	m.keymap.Start.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
