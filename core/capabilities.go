package core

type Capabilities struct {
	AnsiColoring  bool
	HttpService   bool
	ExportMIDI    bool
	ImportMelrose bool
	ReceivingMIDI bool
	SendingMIDI   bool
}
