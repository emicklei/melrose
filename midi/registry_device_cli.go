package midi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emicklei/melrose/notify"
)

func (d *DeviceRegistry) HandleSetting(name string, values []interface{}) error {
	switch name {
	case "midi.in":
		if len(values) != 1 {
			return fmt.Errorf("one argument expected")
		}
		id, ok := values[0].(int)
		if !ok {
			return fmt.Errorf("integer device argument expected")
		}
		_, err := d.Input(id)
		if err != nil {
			return fmt.Errorf("bad input device number: %v", err)
		}
		d.defaultInputID = id
	case "midi.out.channel":
		if len(values) != 2 {
			return fmt.Errorf("two argument expected")
		}
		id, ok := values[0].(int)
		if !ok {
			return fmt.Errorf("integer device argument expected")
		}
		ch, ok := values[1].(int)
		if !ok {
			return fmt.Errorf("integer channel argument expected")
		}
		out, err := d.Output(id)
		if err != nil {
			return fmt.Errorf("bad input device number: %v", err)
		}
		out.defaultChannel = ch
	case "midi.out":
		if len(values) != 1 {
			return fmt.Errorf("one argument expected")
		}
		id, ok := values[0].(int)
		if !ok {
			return fmt.Errorf("integer device argument expected")
		}
		_, err := d.Output(id)
		if err != nil {
			return fmt.Errorf("bad output device number: %v", err)
		}
		d.defaultOutputID = id
	}
	return nil
}

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
