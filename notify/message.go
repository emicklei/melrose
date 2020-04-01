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

func Print(m Message) {
	if m == nil {
		return
	}
	switch m.Type() {
	case NotifyInfo:
		printInfo(m.Message())
	case NotifyWarning:
		printWarning(m.Message())
	case NotifyError:
		printError(m.Message())
	}
}

func printInfo(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;32minfo:\033[0m"}, args...)...)
}

func printError(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;31merror:\033[0m"}, args...)...)
}

func printWarning(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;33mwarning:\033[0m"}, args...)...)
}
