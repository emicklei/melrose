package midi

import "fmt"

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
