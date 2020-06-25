package melrose

type PlayContext struct {
	LoopControl     LoopController
	AudioDevice     AudioDevice
	VariableStorage VariableStorage
}

func (p PlayContext) Control() LoopController    { return p.LoopControl }
func (p PlayContext) Device() AudioDevice        { return p.AudioDevice }
func (p PlayContext) Variables() VariableStorage { return p.VariableStorage }
