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

func NewInfof(format string, args ...any) Message {
	return Notification{message: fmt.Sprintf(format, args...), messageType: NotifyInfo}
}

func NewWarningf(format string, args ...any) Message {
	return Notification{message: fmt.Sprintf(format, args...), messageType: NotifyWarning}
}

func NewError(err error) Message {
	if err == nil {
		return nil
	}
	return NewErrorf("%v", err)
}

func NewErrorf(format string, args ...any) Message {
	return Notification{message: fmt.Sprintf(format, args...), messageType: NotifyError}
}
