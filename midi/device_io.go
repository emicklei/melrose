package midi

import "fmt"

// SendPC = program change
func (d *Device) SendPC(channel, data1 int) error {
	if channel < 1 || channel > 16 {
		return fmt.Errorf("invalid MIDI channel:%d", channel)
	}
	return d.outputStream.WriteShort(int64(0xC0|(channel-1)), int64(data1), 0)
}

// SendCC = control change
func (d *Device) SendCC(channel, data1, data2 int) error {
	if channel < 1 || channel > 16 {
		return fmt.Errorf("invalid MIDI channel:%d", channel)
	}
	return d.outputStream.WriteShort(int64(0xB0|(channel-1)), int64(data1), int64(data2))
}

// if channel < then do not include that information
func (d *Device) SendRaw(status, channel, data1, data2 int) error {
	if channel < 1 {
		return d.outputStream.WriteShort(int64(status), int64(data1), int64(data2))
	}
	if channel > 16 {
		return fmt.Errorf("invalid MIDI channel:%d", channel)
	}
	return d.outputStream.WriteShort(int64(status|(channel-1)), int64(data1), int64(data2))
}
