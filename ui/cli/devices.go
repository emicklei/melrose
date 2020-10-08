package cli

import (
	"log"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/midi"
)

func setupAudio(deviceId string) core.AudioDevice {
	d, err := midi.Open()
	if err != nil {
		log.Fatalln("cannot use audio device:", err)
	}
	return d
}
