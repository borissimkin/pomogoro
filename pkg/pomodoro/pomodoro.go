package pomodoro

import (
	"pomogoro/pkg/session"
	"pomogoro/pkg/settings"
	"sort"
	"time"
)

type Pomodoro struct {
	currentSessionType  session.Type
	settings            *settings.Settings
	sessions            map[session.Type]*session.Session
	previousSessionType session.Type
	completed           map[session.Type]int
}

func (p *Pomodoro) totalWorkSessions() int {
	return p.completed[session.Work]
}

func (p *Pomodoro) currentSession() *session.Session {
	return p.sessions[p.currentSessionType]
}

func (p *Pomodoro) sessionsBeforeLongBreak() int {
	return p.settings.WorkSessionsUntilLongBreak - p.totalWorkSessions()%p.settings.WorkSessionsUntilLongBreak
}

func (p *Pomodoro) getDuration() time.Duration {
	return p.settings.GetDuration(p.currentSessionType)
}

func (p *Pomodoro) getNextSessionType() session.Type {
	if p.previousSessionType != session.Work {
		return session.Work
	}

	if p.settings.WorkSessionsUntilLongBreak <= 0 {
		return session.Break
	}

	if p.completed[session.Work]%p.settings.WorkSessionsUntilLongBreak == 0 {
		return session.LongBreak
	}

	return session.Break
}

func (p *Pomodoro) nextSession() session.Type {
	p.completed[p.currentSessionType]++
	p.previousSessionType = p.currentSessionType

	nextSession := p.getNextSessionType()
	p.currentSessionType = nextSession

	return nextSession
}

func (p *Pomodoro) setSession(session session.Type) {
	p.currentSessionType = session
}

func NewPomodoro(settings *settings.Settings) *Pomodoro {
	return &Pomodoro{
		currentSessionType: session.Work,
		completed:          make(map[session.Type]int),
		settings:           settings,
		sessions: map[session.Type]*session.Session{
			session.Work:      &session.WorkSession,
			session.Break:     &session.BreakSession,
			session.LongBreak: &session.LongBreakSession,
		},
	}
}

func (p *Pomodoro) SliceSessions() []*session.Session {
	sessions := make([]*session.Session, 0, len(p.sessions))

	for _, value := range p.sessions {
		sessions = append(sessions, value)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].SessionType < sessions[j].SessionType
	})

	return sessions
}

func getSessionType(sessionType session.Type) session.Type {
	sessionTypes := session.Types()

	minSessionType := sessionTypes[0]
	maxSessionType := sessionTypes[len(sessionTypes)-1]

	if sessionType < minSessionType {
		return maxSessionType
	}

	if sessionType > maxSessionType {
		return minSessionType
	}

	return sessionType
}
