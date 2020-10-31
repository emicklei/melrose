package core

import "fmt"

type ChannelSelector struct {
	Target Sequenceable
	Number Valueable
}

func NewChannelSelector(target Sequenceable, channel Valueable) ChannelSelector {
	return ChannelSelector{Target: target, Number: channel}
}

func (c ChannelSelector) S() Sequence {
	return c.Target.S()
}

func (c ChannelSelector) Channel() int {
	return Int(c.Number)
}

func (c ChannelSelector) Storex() string {
	if s, ok := c.Target.(Storable); ok {
		return fmt.Sprintf("channel(%v,%s)", c.Number, s.Storex())
	}
	return ""
}

type DeviceSelector struct {
	Target Sequenceable
	ID     Valueable
}

func NewDeviceSelector(target Sequenceable, deviceID Valueable) DeviceSelector {
	return DeviceSelector{Target: target, ID: deviceID}
}

func (d DeviceSelector) S() Sequence {
	return d.Target.S()
}

func (d DeviceSelector) DeviceID() int {
	return Int(d.ID)
}

func (d DeviceSelector) Storex() string {
	if s, ok := d.Target.(Storable); ok {
		return fmt.Sprintf("device(%v,%s)", d.ID, s.Storex())
	}
	return ""
}
