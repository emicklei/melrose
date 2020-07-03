package notify

import "fmt"

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

// Println is to inject a function that can report info,error and warning
var Println = fmt.Println

func printInfo(args ...interface{}) {
	Println(append([]interface{}{"\033[1;32minfo:\033[0m"}, args...)...)
}

func printError(args ...interface{}) {
	Println(append([]interface{}{"\033[1;31merror:\033[0m"}, args...)...)
}

func printWarning(args ...interface{}) {
	Println(append([]interface{}{"\033[1;33mwarning:\033[0m"}, args...)...)
}

func Debugf(format string, args ...interface{}) {
	m := fmt.Sprintf(format, args...)
	Println(append([]interface{}{"\033[1;33mwarning:\033[0m"}, []interface{}{m}...)...)
}
