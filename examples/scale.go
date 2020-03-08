package main

import (
	"time"

	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

var Audio *audio.Device

// go run scale.go
func main() {
	Audio = new(audio.Device)
	Audio.Open()
	Audio.LoadSounds()
	defer Audio.Close()

	cm := C().Scale(Major)

	Audio.Play(cm.S(), 1*time.Second)
}
