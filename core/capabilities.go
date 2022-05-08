package core

type Capabilities struct {
	AnsiColoring  bool
	HTTPService   bool
	ExportMIDI    bool
	ImportMelrose bool
	ReceivingMIDI bool
	SendingMIDI   bool
}

func NewCapabilities() *Capabilities {
	return &Capabilities{
		AnsiColoring:  true,
		HTTPService:   true,
		ExportMIDI:    true,
		ImportMelrose: true,
		ReceivingMIDI: true,
		SendingMIDI:   true,
	}
}
