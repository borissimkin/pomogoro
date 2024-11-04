package settings

import (
	"pomogoro/pkg/session"
	"time"
)

type durations map[session.Type]time.Duration

type Notification struct {
	Sound bool
	Push  bool
}

type AutoStart map[session.Type]bool

type Settings struct {
	WorkSessionsUntilLongBreak int
	Durations                  durations
	ShowProgressBar            bool
	Notification               Notification
	AutoStart                  AutoStart
}

func DefaultSettings() Settings {
	return Settings{
		WorkSessionsUntilLongBreak: 4,
		ShowProgressBar:            true,
		Notification: Notification{
			Sound: true,
			Push:  true,
		},
		AutoStart: AutoStart{
			session.Work:      true,
			session.Break:     true,
			session.LongBreak: true,
		},
		Durations: durations{
			session.Work:      time.Minute * 25,
			session.Break:     time.Minute * 5,
			session.LongBreak: time.Minute * 15,
		},
	}
}

func NewSettings() *Settings {
	old := newStorage().Read()
	if old != nil {
		return old
	}

	s := DefaultSettings()

	return &s
}

func (s *Settings) GetDuration(sessionType session.Type) time.Duration {
	return s.Durations[sessionType]
}
