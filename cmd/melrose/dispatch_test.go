package main

import "testing"

func TestMultiLineEval(t *testing.T) {
	input := `seq("
	C D E C 
	C D E C 
	E F 2G
	E F 2G 
	8G 8A 8G 8F E C 
	8G 8A 8G 8F E C
	2C 2G3 2C
	2C 2G3 2C
	")`
	t.Log(input)
	_, err := eval(input)
	if err != nil {
		t.Error(err)
	}
}
