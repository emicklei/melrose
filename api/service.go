package api

type Service interface {
	CommandInspect(file string, lineEnd int, source string) (interface{}, error)
	CommandPlay(file string, lineEnd int, source string) (interface{}, error)
	CommandStop(file string, lineEnd int, source string) (interface{}, error)
	CommandEvaluate(file string, lineEnd int, source string) (interface{}, error)
	CommandKill() error
	CommandHover(source string) string
	ChangeDefaultDeviceAndChannel(isInput bool, deviceID int, channel int) error
}
