package core

import "sync"

type PlayContext struct {
	LoopControl     LoopController
	AudioDevice     AudioDevice
	VariableStorage VariableStorage
	EnvironmentVars *sync.Map
	CapabilityFlags *Capabilities
}

func (p PlayContext) Control() LoopController     { return p.LoopControl }
func (p PlayContext) Device() AudioDevice         { return p.AudioDevice }
func (p PlayContext) Variables() VariableStorage  { return p.VariableStorage }
func (p PlayContext) Environment() *sync.Map      { return p.EnvironmentVars }
func (p PlayContext) Capabilities() *Capabilities { return p.CapabilityFlags }

type ConditionalPlayContext struct {
	Context
	playCondition Condition
}

func (c ConditionalPlayContext) Condition() Condition { return c.playCondition }

func (p PlayContext) WithCondition(c Condition) Context {
	return ConditionalPlayContext{Context: p, playCondition: c}
}
