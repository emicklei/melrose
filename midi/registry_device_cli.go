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
		return notify.Infof("printing notes enabled:%v", od.echo)
	case "channel":
		if len(args) != 3 {
			return notify.Warningf("missing channel number or device id")
		}
		id, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		nr, err := strconv.Atoi(strings.TrimSpace(args[2]))
		if err != nil {
			return notify.Errorf("bad channel number:%v", err)
		}
		if nr < 1 || nr > 16 {
			return notify.Errorf("bad channel number; must be in [1..16]")
		}
		out, err := d.Output(id)
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		out.defaultChannel = nr
		return notify.Infof("output device id:%d has current MIDI channel:%d", id, nr)
	case "in":
		if len(args) != 2 {
			return notify.Warningf("missing device number")
		}
		nr, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		d.defaultInputID = nr
		return notify.Infof("current input device id:%v", nr)
	case "out":
		if len(args) != 2 {
			return notify.Warningf("missing device number")
		}
		nr, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		d.defaultOutputID = nr
		return notify.Infof("current output device id:%v", nr)
	default:
		return notify.Warningf("unknown device access command: %s", args[0])
	}
}

func (d *DeviceRegistry) printInfo() {
	d.streamRegistry.transport.PrintInfo(d.defaultInputID, d.defaultOutputID)

	od, err := d.Output(d.defaultOutputID)
	if err != nil {
		notify.Print(notify.Error(err))
		return
	}
	fmt.Printf("[midi] channel %d = default MIDI output channel\n", od.defaultChannel)
	fmt.Printf("[midi] echo notes = %v\n", od.echo)
}
