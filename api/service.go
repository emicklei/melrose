package api

import "github.com/emicklei/melrose/core"

type Service interface {
	Context() core.Context
	CommandInspect(file string, lineEnd int, source string) (interface{}, error)
	CommandPlay(file string, lineEnd int, source string) (CommandPlayResponse, error)
	CommandStop(file string, lineEnd int, source string) (interface{}, error)
	CommandEvaluate(file string, lineEnd int, source string) (interface{}, error)
	CommandKill() error
	CommandHover(source string) string
	ChangeDefaultDeviceAndChannel(isInput bool, deviceID int, channel int) error
	CommandMIDISample(file string, lineEnd int, source string) ([]byte, error)
	ListDevices() []core.DeviceDescriptor
}
