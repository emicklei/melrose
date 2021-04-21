// +build !wasm

package transport

import (
	"fmt"
	"log"

	"github.com/emicklei/melrose/notify"
	"gitlab.com/gomidi/rtmididrv/imported/rtmidi"
)

func (t RtmidiTransporter) PrintInfo(inID, outID int) {
	notify.PrintHighlighted("available input:")

	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		log.Fatalln("can't open default MIDI in: ", err)
	}
	defer in.Close()
	ports, err := in.PortCount()
	if err != nil {
		log.Fatalln("can't get number of in ports: ", err)
	}
	for i := 0; i < ports; i++ {
		name, err := in.PortName(i)
		if err != nil {
			name = ""
		}
		fmt.Printf("device %d: %s\n", i, name)
	}
	fmt.Println()

	notify.PrintHighlighted("available output:")
	{
		// Outs
		out, err := rtmidi.NewMIDIOutDefault()
		if err != nil {
			log.Fatalln("can't open default MIDI out: ", err)
		}
		defer out.Close()
		ports, err := out.PortCount()
		if err != nil {
			log.Fatalln("can't get number of out ports: ", err)
		}

		for i := 0; i < ports; i++ {
			name, err := out.PortName(i)
			if err != nil {
				name = ""
			}
			fmt.Printf("device %d: %s\n", i, name)
		}
	}
	fmt.Println()
}
