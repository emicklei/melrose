package main

import (
	"time"

	m "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

var Audio *audio.Device

func main() {
	Audio = new(audio.Device)
	Audio.Open()
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
	go Audio.Play(s, 1*time.Second)

	s2, _ := m.ParseNote("=")
	Audio.Play(s2.Repeated(8).Join(s), 1*time.Second)

	/*
		w := new(sync.WaitGroup)
		for c := 0; c < 2; c++ {
			time.Sleep(time.Duration(c*2) * time.Second)
			go func() {
				w.Add(1)
				Audio.Play(s, 1*time.Second)
				w.Done()
			}()
		}
		w.Wait()
	*/
}
