package melrose

import (
	"fmt"
	//"github.com/emicklei/hopwatch"
)

func debug(v ...interface{}) {
	fmt.Printf("%#v\n", v)
	//hopwatch.CallerOffset(3).Dump(v...).Break()
}
