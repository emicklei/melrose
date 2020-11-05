package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
)

func (r *DeviceRegistry) Listen(deviceID int, who core.NoteListener, startOrStop bool) {
	in, err := r.Input(deviceID)
	if err != nil {
		// TODO
		return
	}
	if startOrStop {
		in.listener.start()
		// wait for pending events to be flushed
		time.Sleep(200 * time.Millisecond)
		in.listener.add(who)
	} else {
		in.listener.remove(who)
		// do not stop the listener such that incoming events are just ignored
	}
}
