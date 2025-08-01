package notify

import (
	"fmt"
	"io"
	"strings"
)

var ansiColorsEnabled = true

func SetANSIColorsEnabled(enabled bool) {
	ansiColorsEnabled = enabled
}

func PrintWelcome(version string) {
	tail := " - program your melodies - " + version + " (help = :h, quit = :q or ctrl+c)"
	if ansiColorsEnabled {
		fmt.Println("\033[1;34mmelrōse\033[0m" + tail)
	} else {
		fmt.Fprint(Console.StandardOut, "melrose"+tail+"\n")
	}
}

func PrintBye() {
	if ansiColorsEnabled {
		fmt.Println("\033[1;34mmelrōse\033[0m" + " sings bye!")
	} else {
		fmt.Fprint(Console.StandardOut, "melrōse sings bye!\n")
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

func Debugf(format string, args ...any) {
	if !IsDebug() {
		return
	}
	// make sure it ends with newline
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(Console.StandardOut, format, args...)
}

func Warnf(format string, args ...any) {
	printWarning(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...any) {
	printInfo(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
	printError(fmt.Sprintf(format, args...))
}

func printInfo(args ...any) {
	fmt.Fprintf(Console.StandardOut, "%s\n", args...)
}

func printError(args ...any) {
	if ansiColorsEnabled {
		Println(append([]any{"\033[1;31merror:\033[0m"}, args...)...)
	} else {
		fmt.Fprintf(Console.StandardError, "ERROR: %s\n", args...)
	}
}

func printWarning(args ...any) {
	if ansiColorsEnabled {
		Println(append([]any{"\033[1;33mwarning:\033[0m"}, args...)...)
	} else {
		fmt.Fprintf(Console.StandardOut, "WARN: %s\n", args...)
	}
}

func PrintKeyValue(b io.Writer, k string, v any) {
	if ansiColorsEnabled {
		fmt.Fprintf(b, "\033[94m%s:\033[0m%v ", k, v)
	} else {
		fmt.Fprintf(b, "%s:%v ", k, v)
	}
}
