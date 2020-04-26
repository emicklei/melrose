package melrose

type PlayContext struct {
	Timeline    *Timeline
	LoopControl LoopController
	AudioDevice AudioDevice
}

var globalPlayContext = &PlayContext{
	Timeline:    NewTimeline(),
	LoopControl: NewBeatmaster(120.0),
	AudioDevice: nil, // set later
}

func Context() *PlayContext {
	return globalPlayContext
}

func (c *PlayContext) SetCurrentDevice(a AudioDevice) {
	c.AudioDevice = a
}
