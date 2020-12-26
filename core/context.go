package core

type PlayContext struct {
	LoopControl     LoopController
	AudioDevice     AudioDevice
	VariableStorage VariableStorage
	EnvironmentVars map[string]string
}

func (p PlayContext) Control() LoopController        { return p.LoopControl }
func (p PlayContext) Device() AudioDevice            { return p.AudioDevice }
func (p PlayContext) Variables() VariableStorage     { return p.VariableStorage }
func (p PlayContext) Environment() map[string]string { return p.EnvironmentVars }

type ConditionalPlayContext struct {
	Context
	playCondition Condition
}

func (c ConditionalPlayContext) Condition() Condition { return c.playCondition }

func (p PlayContext) WithCondition(c Condition) Context {
	return ConditionalPlayContext{Context: p, playCondition: c}
}
