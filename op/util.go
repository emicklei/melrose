package op

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/emicklei/melrose/core"
)

type int2float32 struct {
	at    int
	float float32
}
type int2int struct {
	from int
	to   int
}

// "1 (4 5 6) 2 (4 5 6) 3 (4 5 6) 2 (4 5 6)"
func parseIndices(src string) [][]int {
	ii := [][]int{}
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	container := []int{}
	ingroup := false
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '(':
			if len(container) > 0 {
				ii = append(ii, container)
			}
			container = []int{}
			ingroup = true
		case ')':
			if len(container) > 0 {
				ii = append(ii, container)
			}
			container = []int{}
			ingroup = false
		default:
			i, err := strconv.Atoi(s.TokenText())
			if err == nil {
				if ingroup {
					container = append(container, i)
				} else {
					ii = append(ii, []int{i})
				}
			}
		}
	}
	if len(container) > 0 {
		ii = append(ii, container)
	}
	return ii
}

func replacedAll(target []core.Sequenceable, from, to core.Sequenceable) []core.Sequenceable {
	newTarget := []core.Sequenceable{}
	for _, each := range target {
		if core.IsIdenticalTo(each, from) {
			newTarget = append(newTarget, to)
		} else {
			if other, ok := each.(core.Replaceable); ok {
				newTarget = append(newTarget, other.Replaced(from, to))
			} else {
				newTarget = append(newTarget, each)
			}
		}
	}
	return newTarget
}

// "1 (4 5 6) 2 (4 5 6) 3 (4 5 6) 2 (4 5 6)"
func formatIndices(src [][]int) string {
	var b bytes.Buffer
	for _, each := range src {
		if len(each) == 1 {
			fmt.Fprintf(&b, "%d ", each[0])
		} else {
			fmt.Fprintf(&b, "(")
			for _, other := range each {
				fmt.Fprintf(&b, "%d ", other)
			}
			fmt.Fprintf(&b, ") ")
		}
	}
	return b.String()
}

// 1:-1, 3:-1 ,1:0,2:0,3:0,1:1,2:1
func parseIndexOffsets(s string) (m []int2int) {
	entries := strings.Split(s, ",")
	for _, each := range entries {
		kv := strings.Split(each, ":")
		ik, err := strconv.Atoi(strings.TrimSpace(kv[0]))
		if err != nil {
			continue
		}
		iv, err := strconv.Atoi(strings.TrimSpace(kv[1]))
		if err != nil {
			continue
		}
		m = append(m, int2int{from: ik, to: iv})
	}
	return
}

// 1:1, 2:1.0, 3:0.5, 4:0.01625, 1:2, 1:4, 1:8, 1:16
func parseIndexFloats(s string) (m []int2float32) {
	entries := strings.Split(s, ",")
	for _, each := range entries {
		kv := strings.Split(each, ":")
		ik, err := strconv.Atoi(strings.TrimSpace(kv[0]))
		if err != nil {
			continue
		}
		iv, err := strconv.ParseFloat(strings.TrimSpace(kv[1]), 32)
		if err != nil {
			continue
		}
		m = append(m, int2float32{at: ik, float: float32(iv)})
	}
	return
}
