package main

import (
	"log"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/midi"
)

func setupAudio(deviceId string) melrose.AudioDevice {
	d, err := midi.Open()
	if err != nil {
		log.Fatalln("cannot use audio device:", err)
	}
	return d
}
