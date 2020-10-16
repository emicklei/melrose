package notify

import (
	"fmt"
	"io"
	"os"
)

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
	fmt.Fprintf(c.StandardError, format, args...)
}
