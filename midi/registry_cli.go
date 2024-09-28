package midi

import (
	"fmt"
	"strconv"

	"github.com/emicklei/melrose/notify"
)

func (r *DeviceRegistry) HandleSetting(name string, values []interface{}) error {
	switch name {
	case "echo": // i|o id
		isInput := values[0] == "i"
		id := values[1].(int)
		r.toggleEchoNotesForDevice(isInput, id)
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
	if len(args) == 0 {
		r.printInfoVerbose()
		return nil
	}
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
	if len(args) == 3 && args[0] == "e" {
		if args[1] != "i" && args[1] != "o" {
			return notify.NewErrorf("first parameter is either `i` for input or `o` for output")
		}
		id, err := strconv.Atoi(args[2])
		if err != nil {
			return notify.NewError(err)
		}
		r.HandleSetting("echo", []any{args[1], id})
		return nil
	}
	if len(args) == 1 && args[0] == "r" {
		fmt.Println("Reset MIDI device configuration. Stopping all listeners")
		r.Reset()
		r.Close()
		r.init()
	}
	return notify.NewErrorf("unknown command:%v", args)
}

func (r *DeviceRegistry) printInfo() {
	r.streamRegistry.transport.PrintInfo(r.defaultInputID, r.defaultOutputID)

}
func (r *DeviceRegistry) printInfoVerbose() {
	r.printInfo()

	notify.PrintHighlighted("default settings:")
	deviceIn, err := r.Input(r.defaultInputID)
	if err == nil {
		fmt.Printf(" input device = %d, echo = %v\n", r.defaultInputID, deviceIn.echo)
	} else {
		fmt.Printf(" no input device\n")
	}
	deviceOut, err := r.Output(r.defaultOutputID)
	if err == nil {
		fmt.Printf("output device = %d, channel = %d, echo = %v\n", r.defaultOutputID, deviceOut.defaultChannel, deviceOut.echo)
	} else {
		fmt.Printf(" no output device (restart?)\n")
	}
	fmt.Println()

	notify.PrintHighlighted("how to change:")
	fmt.Println("set('midi.in', <device-id>)              --- change the default MIDI input device id (or use e.g. \":m i 1\")")
	fmt.Println("set('midi.out',<device-id>)              --- change the default MIDI output device id (or use e.g. \":m o 1\")")
	fmt.Println("set('midi.out.channel',<device-id>,<nr>) --- change the default MIDI channel for an output device id")
	fmt.Println(":e i <device-id>                         --- toggle printing the MIDI notes from input device id")
	fmt.Println(":e o <device-id>                         --- toggle printing the MIDI notes to output device id")
}

func (r *DeviceRegistry) toggleEchoNotesForDevice(isInput bool, deviceID int) {
	if isInput {
		in, ok := r.in[deviceID]
		if !ok {
			notify.Errorf("no device found with id:%d", deviceID)
			return
		}
		in.echo = !in.echo
		if in.echo {
			in.listener.Add(DefaultEchoListener)
			in.listener.Start()
		} else {
			in.listener.Remove(DefaultEchoListener)
			// each.listener.Stop()
		}
		notify.Infof("echo input notes from device %d (%s) is enabled: %v", in.id, in.name, in.echo)
	} else {
		out, ok := r.out[deviceID]
		if !ok {
			notify.Errorf("no device found with id:%d", deviceID)
			return
		}
		out.echo = !out.echo
		notify.Infof("echo output notes from device %d (%s) is enabled: %v", out.id, out.name, out.echo)
	}
}
