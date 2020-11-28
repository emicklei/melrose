package notify

import (
	"fmt"
	"strings"
)

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

// func printInfo(args ...interface{}) {
// 	fmt.Fprintf(Console.StandardOut, "INFO: %s\n", args...)
// }

// func printError(args ...interface{}) {
// 	fmt.Fprintf(Console.StandardError, "ERROR: %s\n", args...)
// }

// func printWarning(args ...interface{}) {
// 	fmt.Fprintf(Console.StandardOut, "WARN: %s\n", args...)
// }

func Debugf(format string, args ...interface{}) {
	// make sure it ends with newline
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(Console.StandardOut, format, args...)
}

func printInfo(args ...interface{}) {
	Println(append([]interface{}{"\033[1;32minfo:\033[0m"}, args...)...)
}

func printError(args ...interface{}) {
	Println(append([]interface{}{"\033[1;31merror:\033[0m"}, args...)...)
}

func printWarning(args ...interface{}) {
	Println(append([]interface{}{"\033[1;33mwarning:\033[0m"}, args...)...)
}
