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

func printInfo(args ...interface{}) {
	fmt.Fprintf(Console.StandardOut, "INFO: %s\n", args...)
}

func printError(args ...interface{}) {
	fmt.Fprintf(Console.StandardError, "ERROR: %s\n", args...)
}

func printWarning(args ...interface{}) {
	fmt.Fprintf(Console.StandardOut, "WARN: %s\n", args...)
}

func Debugf(format string, args ...interface{}) {
	// make sure it ends with newline
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(Console.StandardOut, format, args...)
}
