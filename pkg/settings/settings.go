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

type AutoStart struct {
	Break     bool
	Work      bool
	LongBreak bool
}

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
			Break:     true,
			LongBreak: true,
			Work:      true,
		},
		Durations: durations{
			session.Work:      time.Minute * 25,
			session.Break:     time.Minute * 5,
			session.LongBreak: time.Minute * 15,
			//session.Work:      time.Second * 4,
			//session.Break:     time.Second * 3,
			//session.LongBreak: time.Second * 7,
		},
	}
}

func NewSettings() *Settings {
	s := DefaultSettings()

	return &s
}

func (s *Settings) GetDuration(sessionType session.Type) time.Duration {
	return s.Durations[sessionType]
}

// todo: Save to json & load
