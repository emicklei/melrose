package core

import (
	"bytes"
	"fmt"
)

type ChannelSelector struct {
	Target []Sequenceable
	Number Valueable
}

func NewChannelSelector(target []Sequenceable, channel Valueable) ChannelSelector {
	return ChannelSelector{Target: target, Number: channel}
}

func (c ChannelSelector) S() Sequence {
	groups := [][]Note{}
	for _, each := range c.Target {
		groups = append(groups, each.S().Notes...)
	}
	return Sequence{Notes: groups}
}

func (c ChannelSelector) Unwrap() Sequenceable {
	if len(c.Target) == 1 {
		return c.Target[0]
	}
	return SequenceableSlice(c.Target)
}

func (c ChannelSelector) Channel() int {
	return Int(c.Number)
}

func (c ChannelSelector) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "channel(%v", c.Number)
	AppendStorexList(&b, false, c.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

type DeviceSelector struct {
	Target []Sequenceable
	ID     Valueable
}

func NewDeviceSelector(target []Sequenceable, deviceID Valueable) DeviceSelector {
	return DeviceSelector{Target: target, ID: deviceID}
}

func (d DeviceSelector) S() Sequence {
	groups := [][]Note{}
	for _, each := range d.Target {
		groups = append(groups, each.S().Notes...)
	}
	return Sequence{Notes: groups}
}

func (d DeviceSelector) Unwrap() Sequenceable {
	if len(d.Target) == 1 {
		return d.Target[0]
	}
	return SequenceableSlice(d.Target)
}

func (d DeviceSelector) DeviceID() int {
	return Int(d.ID)
}

func (d DeviceSelector) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "device(%v", d.ID)
	AppendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// SequenceableSlice is minimal Join
type SequenceableSlice []Sequenceable

func (s SequenceableSlice) S() Sequence {
	if len(s) == 0 {
		return EmptySequence
	}
	head := s[0].S()
	for i := 1; i < len(s); i++ {
		head = head.SequenceJoin(s[i].S())
	}
	return head
}
