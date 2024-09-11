package core

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/emicklei/melrose/notify"
)

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

func IsIdenticalTo(left, right Sequenceable) bool {
	return reflect.DeepEqual(left, right)
}

func PrintValue(ctx Context, v interface{}) {
	InspectValue(ctx, v)
}

func InspectValue(ctx Context, v interface{}) {
	if v == nil {
		return
	}
	varname := ctx.Variables().NameFor(v)
	i := NewInspect(ctx, varname, v)
	fmt.Fprintf(notify.Console.StandardOut, "%s\n", i.String())
}

func Storex(v interface{}) string {
	if s, ok := v.(Storable); ok {
		return s.Storex()
	}
	if s, ok := v.(string); ok {
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("%v", v)
}

func AppendStorexList(b *bytes.Buffer, isFirstParameter bool, list []Sequenceable) {
	if len(list) == 0 {
		return
	}
	if !isFirstParameter {
		fmt.Fprintf(b, ",")
	}
	for i, each := range list {
		if s, ok := each.(Storable); !ok {
			fmt.Fprintf(b, "nil")
		} else {
			fmt.Fprintf(b, "%s", s.Storex())
		}
		if i < len(list)-1 {
			io.WriteString(b, ",")
		}
	}
}

func UnValue(v Sequenceable) Sequenceable {
	if s, ok := v.(HasValue); ok {
		if seq, ok := s.Value().(Sequenceable); ok {
			return seq
		}
		return EmptySequence
	}
	return v
}

func InList(s Sequenceable) []Sequenceable {
	if s == nil {
		return []Sequenceable{}
	}
	return []Sequenceable{s}
}

func ContainsInt(list []int, value int) bool {
	for _, each := range list {
		if each == value {
			return true
		}
	}
	return false
}

func ReplacedAll(target []Sequenceable, from, to Sequenceable) []Sequenceable {
	newTarget := []Sequenceable{}
	for _, each := range target {
		if IsIdenticalTo(each, from) {
			newTarget = append(newTarget, to)
		} else {
			if other, ok := each.(Replaceable); ok {
				newTarget = append(newTarget, other.Replaced(from, to))
			} else {
				newTarget = append(newTarget, each)
			}
		}
	}
	return newTarget
}
