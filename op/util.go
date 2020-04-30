package op

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emicklei/melrose"
)

func appendStorexList(b *bytes.Buffer, isFirstParameter bool, list []melrose.Sequenceable) {
	if len(list) == 0 {
		return
	}
	if !isFirstParameter {
		fmt.Fprintf(b, ",")
	}
	for i, each := range list {
		if s, ok := each.(melrose.Storable); !ok {
			fmt.Fprintf(b, "nil")
		} else {
			fmt.Fprintf(b, "%s", s.Storex())
		}
		if i < len(list)-1 {
			io.WriteString(b, ",")
		}
	}
}
