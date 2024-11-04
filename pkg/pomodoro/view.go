package pomodoro

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"time"
)

var (
	tabStyles = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			MarginRight(1).
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)
	timerStyles            = lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Bold(true)
	progressBarPausedColor = "#4b4453"
)

func isPause(m *Model) bool {
	return m.keymap.Start.Enabled()
}

func renderProgressBar(m *Model) string {
	color := m.pomodoro.currentSession().BackgroundColor

	if isPause(m) {
		color = progressBarPausedColor
	}

	m.progress.FullColor = color

	return m.progress.ViewAs(getPercent(m))
}

func formatTime(t time.Duration) string {
	return t.Truncate(time.Second).String()
}

func renderTime(m *Model) string {
	var style = timerStyles

	if !m.timer.Running() {
		style = style.Faint(true)
	}

	return style.Render(formatTime(m.timer.Timeout))
}

func getPercent(m *Model) float64 {
	return float64(m.initTime-m.timer.Timeout) / float64(m.initTime)
}

func renderTotalSessions(p *Pomodoro) string {
	return fmt.Sprintf("Total Work sessions: %v", p.totalWorkSessions())
}

func renderSessionsBeforeLongBreak(p *Pomodoro) string {
	return fmt.Sprintf("Sessions left before the long break: %v", p.sessionsBeforeLongBreak())
}

func renderSessionTypes(p *Pomodoro) string {
	s := ""

	for _, item := range p.SliceSessions() {
		cursor := " "

		var style = tabStyles.Background(lipgloss.Color(item.BackgroundColor))

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

func renderBreakLine() string {
	return "\n"
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
