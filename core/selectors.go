package core

import (
	"bytes"
	"fmt"
)

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

func (c ChannelSelector) Unwrap() Sequenceable {
	return c.Target
}

func (c ChannelSelector) Channel() int {
	return Int(c.Number)
}

func (c ChannelSelector) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "channel(%v,%s", c.Number, Storex(c.Target))
	fmt.Fprintf(&b, ")")
	return b.String()
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

func (d DeviceSelector) Unwrap() Sequenceable {
	return d.Target
}

func (d DeviceSelector) DeviceID() int {
	return Int(d.ID)
}

func (d DeviceSelector) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "device(%v,%s", d.ID, Storex(d.Target))
	fmt.Fprintf(&b, ")")
	return b.String()
}
