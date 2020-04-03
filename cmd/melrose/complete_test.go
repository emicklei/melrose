package main

import (
	"testing"
)

func TestCompleteMe(t *testing.T) {
	line := "se"
	head, c, tail := completeMe(line, 0)
	t.Logf("head=%q,completions:%v,tail=%q", head, c, tail)
}

func TestCompleteMe2(t *testing.T) {
	line := ""
	head, c, tail := completeMe(line, 0)
	t.Logf("head=%q,completions:%v,tail=%q", head, c, tail)
}

func TestCompleteMe3(t *testing.T) {
	line := "a = seq"
	head, c, tail := completeMe(line, 7)
	t.Logf("head=%q,completions:%v,tail=%q", head, c, tail)
}
