package pomodoro

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"pomogoro/pkg/app"
	"pomogoro/pkg/notification"
	"pomogoro/pkg/pomodoro/keybinding"
	"pomogoro/pkg/router"
	"pomogoro/pkg/settings"
	"time"
)

const (
	progressBarMaxWidth = 43
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
			return m.router.To(app.SettingsPageName)
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
