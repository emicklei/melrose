package control

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Key struct {
	deviceID int
	channel  int
	note     core.Note
}

func NewKey(deviceID int, channel int, note core.Note) Key {
	return Key{deviceID: deviceID, channel: channel, note: note}
}

// Inspect is part of Inspectable
func (k Key) Inspect(i core.Inspection) {
	i.Properties["device"] = k.deviceID
	i.Properties["channel"] = k.channel
	i.Properties["note"] = k.note
}

func (k Key) DeviceID() int   { return k.deviceID }
func (k Key) Note() core.Note { return k.note }

// Storex is part of core.Storable
func (k Key) Storex() string {
	return fmt.Sprintf("key(device(%d,channel(%d,%s)))", k.deviceID, k.channel, k.note.Storex())
}
