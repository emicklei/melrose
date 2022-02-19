package midi

import (
	"fmt"
	"strconv"

	"github.com/emicklei/melrose/notify"
)

func (d *DeviceRegistry) HandleSetting(name string, values []interface{}) error {
	switch name {
	case "echo":
		if len(values) != 1 {
			return fmt.Errorf("one argument expected")
		}
		enable, ok := values[0].(bool)
		if !ok {
			return fmt.Errorf("boolean device argument expected, got %T", values[0])
		}
		od, _ := d.Output(d.defaultOutputID)
		od.echo = enable
		notify.Infof("echo notes is enabled: %v", enable)
	case "echo.toggle":
		if len(values) != 0 {
			return fmt.Errorf("no argument expected")
		}
		od, _ := d.Output(d.defaultOutputID)
		od.echo = !od.echo
		notify.Infof("echo notes is enabled: %v", od.echo)
	case "midi.in":
		if len(values) != 1 {
			return fmt.Errorf("one argument expected")
		}
		id, ok := values[0].(int)
		if !ok {
			return fmt.Errorf("integer device argument expected, got %T", values[0])
		}
		_, err := d.Input(id)
		if err != nil {
			return fmt.Errorf("bad input device number: %v", err)
		}
		d.defaultInputID = id
		notify.Infof("Set default MIDI input device id: %d", id)
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
		notify.Infof("Set default MIDI output device id: %d with default channel: %d", id, ch)
	case "midi.out":
		if len(values) != 1 {
			return fmt.Errorf("one argument expected")
		}
		id, ok := values[0].(int)
		if !ok {
			return fmt.Errorf("integer device argument expected")
		}
		out, err := d.Output(id)
		if err != nil {
			return fmt.Errorf("bad output device number: %v", err)
		}
		d.defaultOutputID = id
		notify.Infof("Set default MIDI output device id: %d with default channel: %d", id, out.defaultChannel)
	default:
		return fmt.Errorf("unknown setting:%s", name)
	}
	return nil
}

// Command is part of melrose.AudioDevice
// TODO obsolete?
func (d *DeviceRegistry) Command(args []string) notify.Message {
	fmt.Println(args)
	if len(args) == 2 && args[0] == "o" {
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.NewError(err)
		}
		if err := d.HandleSetting("midi.out", []interface{}{id}); err != nil {
			return notify.NewError(err)
		}
		return nil
	}
	if len(args) == 2 && args[0] == "i" {
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.NewError(err)
		}
		if err := d.HandleSetting("midi.in", []interface{}{id}); err != nil {
			return notify.NewError(err)
		}
		return nil
	}
	if len(args) == 1 && args[0] == "e" {
		d.HandleSetting("echo.toggle", []interface{}{})
		return nil
	}
	if len(args) == 1 && args[0] == "r" {
		fmt.Println("Reset MIDI device configuration. Stopping all listeners")
		d.Reset()
		d.Close()
		d.init()
	}
	d.printInfo()
	return nil
}

func (d *DeviceRegistry) printInfo() {
	d.streamRegistry.transport.PrintInfo(d.defaultInputID, d.defaultOutputID)

	notify.PrintHighlighted("current defaults:")
	_, err := d.Input(d.defaultInputID)
	if err == nil {
		fmt.Printf(" input device = %d\n", d.defaultInputID)
	} else {
		fmt.Printf(" no input device\n")
	}
	od, err := d.Output(d.defaultOutputID)
	if err == nil {
		fmt.Printf("output device = %d, channel = %d\n", d.defaultOutputID, od.defaultChannel)
		fmt.Printf("   echo notes = %v\n", od.echo)
	} else {
		fmt.Printf(" no output device (restart?)\n")
	}

	fmt.Println()

	notify.PrintHighlighted("change:")
	fmt.Println("set('midi.in',<device-id>)               --- change the default MIDI input device id (or e.g. \":m i 1\")")
	fmt.Println("set('midi.out',<device-id>)              --- change the default MIDI output device id (or e.g. \":m o 1\")")
	fmt.Println("set('midi.out.channel',<device-id>,<nr>) --- change the default MIDI channel for an output device id")
	fmt.Println("set('echo.toggle')                       --- toggle printing the notes (or \":m e\" )")
	fmt.Println("set('echo',true)                         --- true = print the notes")
}
