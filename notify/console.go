package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var debugEnabled = false

func IsDebug() bool {
	return debugEnabled
}

func ToggleDebug() bool {
	debugEnabled = !debugEnabled
	return debugEnabled
}

var Console = ConsoleWriter{
	DeviceIn:      os.Stdout,
	DeviceOut:     os.Stdout,
	StandardOut:   os.Stdout,
	StandardError: os.Stderr,
}

type ConsoleWriter struct {
	DeviceIn      io.Writer
	DeviceOut     io.Writer
	StandardOut   io.Writer
	StandardError io.Writer
}

func (c ConsoleWriter) Errorf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(c.StandardError, format, args...)
}

func (c ConsoleWriter) Warnf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(c.StandardError, format, args...)
}
