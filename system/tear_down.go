package system

import (
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

func TearDown(ctx core.Context) error {
	dsl.StopAllPlayables(ctx)
	ctx.Control().Reset()
	ctx.Device().Close()
	notify.PrintBye()
	return nil
}
