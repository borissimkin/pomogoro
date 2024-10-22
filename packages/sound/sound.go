package sound

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Player struct {
	player      *oto.Player
	initialized bool
}

func NewSoundPlayer() *Player {
	return &Player{
		initialized: false,
	}
}

func (s *Player) InitSoundContext() {
	fileBytes, err := os.ReadFile("./assets/ring.mp3")
	if err != nil {
		return
	}

	fileBytesReader := bytes.NewReader(fileBytes)

	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	op := &oto.NewContextOptions{}

	op.SampleRate = 44100

	op.ChannelCount = 2

	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return
	}
	<-readyChan

	player := otoCtx.NewPlayer(decodedMp3)

	s.player = player
	s.initialized = true
}

func (s *Player) Play() {
	if !s.initialized {
		return
	}
	s.player.Play()

	for s.player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	_, err := s.player.Seek(0, io.SeekStart)
	if err != nil {
		s.initialized = false
	}
}
