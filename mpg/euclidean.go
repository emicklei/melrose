package mpg

import "github.com/emicklei/melrose/core"

type Euclidean struct {
	Steps    core.HasValue
	Pulses   core.HasValue
	Rotation core.HasValue
	// playback
	Rate       core.HasValue
	NoteLength core.HasValue
	Channel    int
	Pitch      core.HasValue
	Velocity   core.HasValue
}
