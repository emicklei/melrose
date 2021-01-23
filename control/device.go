package control

type ChannelOnDevice struct {
	input   bool
	name    string
	id      int
	channel int
}

func NewChannelOnDevice(input bool,
	name string,
	id int,
	channel int) ChannelOnDevice {
	return ChannelOnDevice{
		input:   input,
		name:    name,
		id:      id,
		channel: channel,
	}
}
