package main

import "time"

type Pomodoro struct {
	currentSessionType  SessionType
	settings            *Settings
	sessionSettings     map[SessionType]*SessionSettings
	previousSessionType SessionType
	completed           map[SessionType]int
}

func (p *Pomodoro) totalWorkSessions() int {
	return p.completed[workSession]
}

func (p *Pomodoro) currentSessionSettings() *SessionSettings {
	return p.sessionSettings[p.currentSessionType]
}

func (p *Pomodoro) sessionsBeforeLongBreak() int {
	return p.settings.workSessionsUntilLongBreak - p.totalWorkSessions()%p.settings.workSessionsUntilLongBreak
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

func (p *Pomodoro) nextSession() SessionType {
	p.completed[p.currentSessionType]++
	p.previousSessionType = p.currentSessionType

	nextSession := p.getNextSessionType()
	p.currentSessionType = nextSession

	return nextSession
}

func (p *Pomodoro) setSession(session SessionType) {
	p.currentSessionType = session
}
