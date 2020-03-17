package main

import (
	"fmt"
	"testing"

	"github.com/emicklei/melrose"
)

func TestAvailableNoteMethods(t *testing.T) {
	n := melrose.C()
	list := availableMethodNames(n, "")
	fmt.Println(list)
}
