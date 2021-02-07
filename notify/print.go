package notify

import (
	"fmt"
	"strings"
)

var ansiColorsEnabled = true

func PrintWelcome() {
	if ansiColorsEnabled {
		fmt.Println("\033[1;34mmelrōse\033[0m" + " - program your melodies")
	} else {
		fmt.Fprintf(Console.StandardOut, "melrose - program your melodies\n")
	}
}

func PrintBye() {
	if ansiColorsEnabled {
		fmt.Println("\033[1;34mmelrose\033[0m" + " sings bye!")
	} else {
		fmt.Fprintf(Console.StandardOut, "melrose sings bye!\n")
	}
}

func Prompt() string {
	if ansiColorsEnabled {
		return "𝄞 "
	}
	return "# "
}

func PrintHighlighted(what string) {
	if ansiColorsEnabled {
		fmt.Println("\033[1;33m" + what + "\033[0m")
	} else {
		fmt.Println(what)
	}
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

// Println is to inject a function that can report info,error and warning
var Println = fmt.Println

func Debugf(format string, args ...interface{}) {
	// make sure it ends with newline
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(Console.StandardOut, format, args...)
}

func Warnf(format string, args ...interface{}) {
	printWarning(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...interface{}) {
	printInfo(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	printError(fmt.Sprintf(format, args...))
}

func printInfo(args ...interface{}) {
	if ansiColorsEnabled {
		Println(append([]interface{}{"\033[1;32minfo:\033[0m"}, args...)...)
	} else {
		fmt.Fprintf(Console.StandardOut, "INFO: %s\n", args...)
	}
}

func printError(args ...interface{}) {
	if ansiColorsEnabled {
		Println(append([]interface{}{"\033[1;31merror:\033[0m"}, args...)...)
	} else {
		fmt.Fprintf(Console.StandardError, "ERROR: %s\n", args...)
	}
}

func printWarning(args ...interface{}) {
	if ansiColorsEnabled {
		Println(append([]interface{}{"\033[1;33mwarning:\033[0m"}, args...)...)
	} else {
		fmt.Fprintf(Console.StandardOut, "WARN: %s\n", args...)
	}
}
