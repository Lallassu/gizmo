package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
	"strings"
	"time"
)

type sound struct {
	sounds map[string]beep.StreamSeekCloser
}

func (s *sound) create() {
	s.sounds = make(map[string]beep.StreamSeekCloser)

	for _, file := range gSoundFiles {
		f, _ := os.Open(fmt.Sprintf("assets/sounds/%v", file))
		Debug("Loading sound", file)
		if strings.Contains(file, "mp3") {
			name := strings.TrimSuffix(file, ".mp3")
			stream, format, err := mp3.Decode(f)
			if err != nil {
				panic("Failed to load sound")
			}
			s.sounds[name] = stream
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		} else if strings.Contains(file, "wav") {
			name := strings.TrimSuffix(file, ".wav")
			stream, format, err := wav.Decode(f)
			if err != nil {
				panic("Failed to load sound")
			}
			s.sounds[name] = stream
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		}
	}
}

func (s *sound) play(name string) {
	speaker.Play(s.sounds[name])
}
