package settings

import (
	"pomogoro/pkg/session"
	"time"
)

type durations map[session.Type]time.Duration

type Settings struct {
	WorkSessionsUntilLongBreak int
	Durations                  durations
}

func NewSettings() *Settings {
	return &Settings{
		WorkSessionsUntilLongBreak: 4,
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

func (s *Settings) GetDuration(sessionType session.Type) time.Duration {
	return s.Durations[sessionType]
}
