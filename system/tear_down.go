package system

import (
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

func TearDown(ctx core.Context) error {
	// Stop timing/scheduling first to minimize new events during shutdown.
	dsl.StopAllPlayables(ctx)
	ctx.Control().Stop()
	// Force immediate silence before closing MIDI streams.
	ctx.Device().Reset()
	if err := ctx.Device().Close(); err != nil {
		return err
	}
	notify.PrintBye()
	return nil
}
