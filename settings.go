package main

import (
	"sort"
	"time"
)

type SessionType int

const (
	workSession      SessionType = 1
	breakSession     SessionType = 2
	longBreakSession SessionType = 3
)

func getSessionType(sessionType SessionType) SessionType {
	sessionTypes := sliceSessionTypes()

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

func sliceSessionTypes() []SessionType {
	sessionTypes := []SessionType{workSession, breakSession, longBreakSession}

	sort.Slice(sessionTypes, func(i, j int) bool {
		return i < j
	})

	return sessionTypes
}

var workSessionSettings = SessionSettings{
	sessionType:     workSession,
	title:           "Pomodoro",
	backgroundColor: "#ba4949",
	notification: &SessionNotifySettings{
		title:   "Work",
		message: "It’s time to focus and make some progress!",
	},
}

var breakSessionSettings = SessionSettings{
	sessionType:     breakSession,
	title:           "Short Break",
	backgroundColor: "#38858a",

	notification: &SessionNotifySettings{
		title:   "Break",
		message: "Take a short break to recharge and reset.",
	},
}

var longBreakSessionSettings = SessionSettings{
	sessionType:     longBreakSession,
	title:           "Long Break",
	backgroundColor: "#397097",
	notification: &SessionNotifySettings{
		title:   "Rest",
		message: "Enjoy a longer break to fully unwind and refresh.",
	},
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
			workSession:  time.Second * 4,
			breakSession: time.Second * 3,
			//longBreakSession: time.Second * 7,
		},
	}
}

type SessionNotifySettings struct {
	title   string
	message string
}

type SessionSettings struct {
	sessionType     SessionType
	title           string
	emoji           int
	backgroundColor string
	notification    *SessionNotifySettings
}
