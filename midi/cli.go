package midi

import (
	"fmt"
	"strconv"

	"github.com/emicklei/melrose/notify"
)

func (r *DeviceRegistry) HandleSetting(name string, values []interface{}) error {
	switch name {
	case "echo":
		if len(values) != 1 {
			return fmt.Errorf("one argument expected")
		}
		enable, ok := values[0].(bool)
		if !ok {
			return fmt.Errorf("boolean device argument expected, got %T", values[0])
		}
		od, _ := r.Output(r.defaultOutputID)
		od.echo = enable
		notify.Infof("echo notes is enabled: %v", enable)
	case "echo.toggle":
		if len(values) != 0 {
			return fmt.Errorf("no argument expected")
		}
		// input
		id, err := r.Input(r.defaultInputID)
		if err == nil {
			id.echo = !id.echo
			if id.echo {
				id.listener.Add(DefaultEchoListener)
				id.listener.Start()
			} else {
				id.listener.Remove(DefaultEchoListener)
				// id.listener.Stop()
			}
			notify.Infof("echo input notes from device %d is enabled: %v", id.id, id.echo)
		} else {
			notify.Infof("echo input notes is disabled ; no input device")
		}
		// output
		od, err := r.Output(r.defaultOutputID)
		if err == nil {
			od.echo = !od.echo
			notify.Infof("echo output notes from device %d is enabled: %v", od.id, od.echo)
		} else {
			notify.Infof("echo output notes is disabled ; no output device")
		}
	case "midi.in":
		if len(values) != 1 {
			return fmt.Errorf("one argument expected")
		}
		id, ok := values[0].(int)
		if !ok {
			return fmt.Errorf("integer device argument expected, got %T", values[0])
		}
		_, err := r.Input(id)
		if err != nil {
			return fmt.Errorf("bad input device number: %v", err)
		}
		r.defaultInputID = id
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
		out, err := r.Output(id)
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
		out, err := r.Output(id)
		if err != nil {
			return fmt.Errorf("bad output device number: %v", err)
		}
		r.defaultOutputID = id
		notify.Infof("Set default MIDI output device id: %d with default channel: %d", id, out.defaultChannel)
	default:
		return fmt.Errorf("unknown setting:%s", name)
	}
	return nil
}

// Command is part of melrose.AudioDevice
func (r *DeviceRegistry) Command(args []string) notify.Message {
	if len(args) == 2 && args[0] == "o" {
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.NewError(err)
		}
		if err := r.HandleSetting("midi.out", []interface{}{id}); err != nil {
			return notify.NewError(err)
		}
		return nil
	}
	if len(args) == 2 && args[0] == "i" {
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.NewError(err)
		}
		if err := r.HandleSetting("midi.in", []interface{}{id}); err != nil {
			return notify.NewError(err)
		}
		return nil
	}
	if len(args) == 1 && args[0] == "e" {
		r.HandleSetting("echo.toggle", []interface{}{})
		return nil
	}
	if len(args) == 1 && args[0] == "r" {
		fmt.Println("Reset MIDI device configuration. Stopping all listeners")
		r.Reset()
		r.Close()
		r.init()
	}
	r.printInfo()
	return nil
}

func (r *DeviceRegistry) printInfo() {
	r.streamRegistry.transport.PrintInfo(r.defaultInputID, r.defaultOutputID)

	notify.PrintHighlighted("current defaults:")
	_, err := r.Input(r.defaultInputID)
	if err == nil {
		fmt.Printf(" input device = %d\n", r.defaultInputID)
	} else {
		fmt.Printf(" no input device\n")
	}
	od, err := r.Output(r.defaultOutputID)
	if err == nil {
		fmt.Printf("output device = %d, channel = %d\n", r.defaultOutputID, od.defaultChannel)
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
