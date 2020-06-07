package notify

import "fmt"

const (
	NotifyInfo = iota
	NotifyWarning
	NotifyError
)

type Message interface {
	Message() string
	Type() int
}

type Notification struct {
	message     string
	messageType int
}

func (n Notification) Message() string { return n.message }
func (n Notification) Type() int       { return n.messageType }

func Infof(format string, args ...interface{}) Message {
	return Notification{message: fmt.Sprintf(format, args...), messageType: NotifyInfo}
}

func Warningf(format string, args ...interface{}) Message {
	return Notification{message: fmt.Sprintf(format, args...), messageType: NotifyWarning}
}

func Error(err error) Message {
	if err == nil {
		return nil
	}
	return Errorf("%v", err)
}

func Errorf(format string, args ...interface{}) Message {
	return Notification{message: fmt.Sprintf(format, args...), messageType: NotifyError}
}
