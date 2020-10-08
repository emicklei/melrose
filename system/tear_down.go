package system

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

func TearDown(ctx core.Context) error {
	ctx.Control().Reset()
	ctx.Device().Close()
	fmt.Println("\033[1;34mmelrose\033[0m" + " sings bye!")
	return nil
}
