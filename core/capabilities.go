package core

type Capabilities struct {
	AnsiColoring  bool
	HttpService   bool
	ExportMIDI    bool
	ImportMelrose bool
	ReceivingMIDI bool
	SendingMIDI   bool
}

func NewCapabilities() *Capabilities {
	return &Capabilities{
		AnsiColoring:  true,
		HttpService:   true,
		ExportMIDI:    true,
		ImportMelrose: true,
		ReceivingMIDI: true,
		SendingMIDI:   true,
	}
}
