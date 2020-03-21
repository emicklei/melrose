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

	s, _ := m.ParseSequence(`
	C D E C 
	C D E C 
	E F 2G
	E F 2G 
	8G 8A 8G 8F E C 
	8G 8A 8G 8F E C
	2C 2G3 2C
	2C 2G3 2C`)
	go Audio.Play(s)

	s2, _ := m.ParseNote("=")
	Audio.Play(s2.Repeated(8).Join(s).S())
}
