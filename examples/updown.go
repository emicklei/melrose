package main

import (
	m "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audiolib"
)

func main() {
	Audio, _ := audiolib.Open()
	Audio.BeatsPerMinute(160)
	Audio.LoadSounds()
	defer Audio.Close()

	left := newStreamer(Audio)
	right := newStreamer(Audio)

	note, _ := m.ParseNote("C1")
	for i := 0; i < 40; i++ {
		left.put(note.Major(i))
		right.put(note.Major(39 - i))
	}
}

type streamer struct {
	notes chan m.Note
}

func (s *streamer) put(n m.Note) {
	s.notes <- n
}

func newStreamer(a *audiolib.Device) *streamer {
	s := new(streamer)
	s.notes = make(chan m.Note)
	go func() {
		for {
			a.Play((<-s.notes).S())
		}
	}()
	return s
}
