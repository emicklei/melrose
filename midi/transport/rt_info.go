//go:build !wasm
// +build !wasm

package transport

import (
	"fmt"
	"log"

	"github.com/emicklei/melrose/notify"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

func (t RtmidiTransporter) PrintInfo(defaultInID, defaultOutID int) {
	notify.PrintHighlighted("available inputs:")

	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		log.Fatalln("can't open default MIDI in: ", err)
	}
	defer in.Close()
	ports, err := in.PortCount()
	if err != nil {
		log.Fatalln("can't get number of input ports: ", err)
	}
	for i := 0; i < ports; i++ {
		name, err := in.PortName(i)
		if err != nil {
			name = ""
		}
		isCurrent := ""
		if i == defaultInID {
			isCurrent = " (default)"
		}
		fmt.Printf(" set('midi.in',%d) --- set MIDI input from %s%s\n", i, name, isCurrent)
	}
	fmt.Println()

	notify.PrintHighlighted("available outputs:")
	{
		// Outs
		out, err := rtmidi.NewMIDIOutDefault()
		if err != nil {
			log.Fatalln("can't open default MIDI out: ", err)
		}
		defer out.Close()
		ports, err := out.PortCount()
		if err != nil {
			log.Fatalln("can't get number of output ports: ", err)
		}

		for i := 0; i < ports; i++ {
			name, err := out.PortName(i)
			if err != nil {
				name = ""
			}
			isCurrent := ""
			if i == defaultInID {
				isCurrent = " (default)"
			}

			fmt.Printf("set('midi.out',%d) --- set MIDI output to%s%s\n", i, name, isCurrent)
		}
	}
	fmt.Println()
}
