package midi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emicklei/melrose/notify"
)

// Command is part of melrose.AudioDevice
func (d *DeviceRegistry) Command(args []string) notify.Message {
	if len(args) == 0 {
		d.printInfo()
		return nil
	}
	switch args[0] {
	case "echo":
		od, _ := d.Output(d.defaultOutputID)
		od.echo = !od.echo
		return notify.NewInfof("printing notes enabled: %v", od.echo)
	case "channel":
		if len(args) != 3 {
			return notify.NewWarningf("missing channel number or device id")
		}
		id, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			return notify.NewErrorf("bad device number: %v", err)
		}
		nr, err := strconv.Atoi(strings.TrimSpace(args[2]))
		if err != nil {
			return notify.NewErrorf("bad channel number: %v", err)
		}
		if nr < 1 || nr > 16 {
			return notify.NewErrorf("bad channel number; must be in [1..16]")
		}
		out, err := d.Output(id)
		if err != nil {
			return notify.NewErrorf("bad device number: %v", err)
		}
		out.defaultChannel = nr
		return notify.NewInfof("output device id: %d has current MIDI channel: %d", id, nr)
	case "in":
		if len(args) != 2 {
			return notify.NewWarningf("missing device number")
		}
		nr, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			return notify.NewErrorf("bad device number: %v", err)
		}
		d.defaultInputID = nr
		return notify.NewInfof("current input device id: %d", nr)
	case "out":
		if len(args) != 2 {
			return notify.NewWarningf("missing device number")
		}
		nr, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			return notify.NewErrorf("bad device number:%v", err)
		}
		d.defaultOutputID = nr
		return notify.NewInfof("current output device id: %d", nr)
	default:
		return notify.NewWarningf("unknown device access command: %s", args[0])
	}
}

func (d *DeviceRegistry) printInfo() {
	d.streamRegistry.transport.PrintInfo(d.defaultInputID, d.defaultOutputID)

	notify.PrintHighlighted("current:")
	od, err := d.Output(d.defaultOutputID)
	if err == nil {
		fmt.Printf("input device %d, channel %d\n", d.defaultOutputID, od.defaultChannel)
	}

	id, err := d.Output(d.defaultInputID)
	if err == nil {
		fmt.Printf("output device %d, channel %d\n", d.defaultInputID, id.defaultChannel)
	}
	fmt.Printf("echo notes = %v\n", od.echo)
}
