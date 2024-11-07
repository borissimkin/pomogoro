package session

import (
	"github.com/borissimkin/pomogoro/pkg/notification"
	"sort"
)

type Type int

type Session struct {
	SessionType     Type
	Title           string
	BackgroundColor string
	NotifyParams    *notification.NotifyParams
}

const (
	Work      Type = 1
	Break     Type = 2
	LongBreak Type = 3
)

var WorkSession = Session{
	SessionType:     Work,
	Title:           "Pomodoro",
	BackgroundColor: "#ba4949",
	NotifyParams: &notification.NotifyParams{
		Title:   "Work",
		Message: "It’s time to focus and make some progress!",
	},
}

var BreakSession = Session{
	SessionType:     Break,
	Title:           "Short Break",
	BackgroundColor: "#38858a",
	NotifyParams: &notification.NotifyParams{
		Title:   "Short Break",
		Message: "Take a short break to recharge and reset.",
	},
}

var LongBreakSession = Session{
	SessionType:     LongBreak,
	Title:           "Long Break",
	BackgroundColor: "#397097",
	NotifyParams: &notification.NotifyParams{
		Title:   "Long Break",
		Message: "Enjoy a longer break to fully unwind and refresh.",
	},
}

func Types() []Type {
	sessionTypes := []Type{Work, Break, LongBreak}

	sort.Slice(sessionTypes, func(i, j int) bool {
		return i < j
	})

	return sessionTypes
}
