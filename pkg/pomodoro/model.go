package pomodoro

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pomogoro/pkg/notification"
	"pomogoro/pkg/pomodoro/keybinding"
	"pomogoro/pkg/router"
	"pomogoro/pkg/session"
	"pomogoro/pkg/settings"
	"sort"
	"strings"
	"time"
)

type Model struct {
	progress    progress.Model
	timer       timer.Model
	initTime    time.Duration
	soundPlayer *notification.Player
	keymap      keybinding.KeyMap
	help        help.Model
	pomodoro    *Pomodoro
	router      *router.Router
}

func (m *Model) initPomodoro() {
	m.pomodoro = NewPomodoro(settings.NewSettings())
}

func (m *Model) Init() tea.Cmd {
	m.initPomodoro()
	setTime(m, m.pomodoro.getDuration())
	return tea.Batch(tea.ClearScreen, m.timer.Init())
}

func getIncreasedTime(timeout time.Duration, minutes time.Duration) time.Duration {
	return timeout + minutes*time.Minute
}

func getDecreasedTime(timeout time.Duration, minutes time.Duration) time.Duration {
	return timeout - minutes*time.Minute
}

func setTime(m *Model, duration time.Duration) {
	m.timer.Timeout = duration
	m.initTime = duration
}

const (
	progressBarMaxWidth = 43
)

// todo: прееделать на ресиверы

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		notifyParams := m.pomodoro.sessions[nextSession].NotifyParams
		if m.pomodoro.settings.Notification.Push {
			notification.Notify(notifyParams.Title, notifyParams.Message)
		}
		if m.pomodoro.settings.Notification.Sound {
			m.soundPlayer.Play()
		}
		setTime(m, m.pomodoro.getDuration())

		if !m.pomodoro.settings.AutoStart[nextSession] {
			return m, m.timer.Stop()
		}

		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Settings):
			return m.router.To("settings")
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Reset):
			setTime(m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Start, m.keymap.Stop):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.keymap.Next):
			m.pomodoro.nextSession()
			setTime(m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Right):
			m.pomodoro.setSession(getSessionType(m.pomodoro.currentSessionType + 1))
			setTime(m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Left):
			m.pomodoro.setSession(getSessionType(m.pomodoro.currentSessionType - 1))
			setTime(m, m.pomodoro.getDuration())
		case key.Matches(msg, m.keymap.Up):
			newTimeout := getIncreasedTime(m.timer.Timeout, keybinding.DefaultStepMinutes)
			m.initTime += keybinding.DefaultStepMinutes * time.Minute
			m.timer.Timeout = newTimeout
		case key.Matches(msg, m.keymap.Down):
			newTimeout := getDecreasedTime(m.timer.Timeout, keybinding.DefaultStepMinutes)
			m.initTime -= keybinding.DefaultStepMinutes * time.Minute
			if newTimeout < 0 {
				m.pomodoro.nextSession()
				setTime(m, m.pomodoro.getDuration())
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

	// todo: refactor (use SliceSessionTypes?)
	sessions := make([]*session.Session, 0, len(p.sessions))

	for _, value := range p.sessions {
		sessions = append(sessions, value)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].SessionType < sessions[j].SessionType
	})

	for _, item := range sessions {
		cursor := " "

		var style = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			MarginRight(1).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color(item.BackgroundColor)).
			Padding(0, 1)

		if item.SessionType != p.currentSessionType {
			style = style.
				Faint(true)
		} else {
			cursor = "*"
		}

		s += style.Render(fmt.Sprintf("%s %s", cursor, item.Title))
	}

	return s
}

// todo: pass count
func renderBreakLine() string {
	return "\n"
}

// todo: TRUNCATE на duration?
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

func renderTime(m *Model) string {
	var style = lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Bold(true)

	if !m.timer.Running() {
		style = style.Faint(true)
	}

	// todo: remove milliseconds cause interval with time.Second has bug
	return style.Render(removeMilliseconds(m.timer.View()))
}

func getPercent(m *Model) float64 {
	return float64(m.initTime-m.timer.Timeout) / float64(m.initTime)
}

func isPause(m *Model) bool {
	return m.keymap.Start.Enabled()
}

func renderProgressBar(m *Model) string {
	color := m.pomodoro.currentSession().BackgroundColor

	if isPause(m) {
		color = "#4b4453" // todo
	}

	m.progress.FullColor = color

	return m.progress.ViewAs(getPercent(m))
}

func (m *Model) View() string {
	s := renderSessionTypes(m.pomodoro)

	s += renderBreakLine()
	s += renderBreakLine()

	s += renderTime(m)

	s += renderBreakLine()

	if m.pomodoro.settings.ShowProgressBar {
		s += renderProgressBar(m)
		s += renderBreakLine()
	}

	s += renderBreakLine()
	s += renderTotalSessions(m.pomodoro)
	s += renderBreakLine()

	if m.pomodoro.settings.WorkSessionsUntilLongBreak > 0 {
		s += renderSessionsBeforeLongBreak(m.pomodoro)
		s += renderBreakLine()
	}

	s += m.help.View(m.keymap)

	return s
}

func NewModel(r *router.Router) *Model {
	soundPlayer := notification.NewSoundPlayer()
	soundPlayer.InitSoundContext()

	p := NewPomodoro(settings.NewSettings())

	initTime := p.getDuration()

	model := &Model{
		progress:    progress.New(progress.WithSolidFill(p.currentSession().BackgroundColor), progress.WithoutPercentage()),
		timer:       timer.NewWithInterval(initTime, time.Millisecond),
		initTime:    initTime,
		pomodoro:    p,
		soundPlayer: soundPlayer,
		keymap:      keybinding.InitKeys(),
		help:        help.New(),
		router:      r,
	}
	model.keymap.Start.SetEnabled(false)

	return model
}
