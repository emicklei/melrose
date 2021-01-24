package dsl

import (
	"fmt"
	"log"
	"strings"

	"github.com/emicklei/melrose/core"
)

func IsCompatibleSyntax(s string) bool {
	if len(s) == 0 {
		// ignore syntax ; you are on your own
		return true
	}
	mm := strings.Split(SyntaxVersion, ".")
	ss := strings.Split(s, ".")
	return mm[0] == ss[0] && ss[1] <= mm[1]
}

type Function struct {
	Title         string
	Description   string
	Prefix        string // for autocomplete
	Alias         string // short notation
	Template      string // for autocomplete in VSC
	Samples       string // for doc generation
	ControlsAudio bool
	Tags          string // space separated
	IsCore        bool   // creates a core musical object
	IsComposer    bool   // can decorate a musical object or other decorations
	Func          interface{}
}

func (f Function) HumanizedTemplate() string {
	r := strings.NewReplacer(
		"${1:", "",
		"${2:", "",
		"${3:", "",
		"${4:", "",
		"}", "")
	return r.Replace(f.Template)
}

func registerFunction(m map[string]Function, k string, f Function) {
	if dup, ok := m[k]; ok {
		log.Fatal("duplicate function key detected:", dup)
	}
	if len(f.Alias) > 0 {
		if dup, ok := m[f.Alias]; ok {
			log.Fatal("duplicate function alias key detected:", dup)
		}
	}
	m[k] = f
	if len(f.Alias) > 0 {
		// modify title
		f.Title = fmt.Sprintf("%s [%s]", f.Title, f.Alias)
		m[f.Alias] = f
	}
}

func getSequenceable(v interface{}) (core.Sequenceable, bool) {
	if s, ok := v.(core.Sequenceable); ok {
		return s, ok
	}
	return nil, false
}

func getSequenceableList(m ...interface{}) (list []core.Sequenceable, ok bool) {
	ok = true
	for _, each := range m {
		if s, ok := getSequenceable(each); ok {
			list = append(list, s)
		} else {
			return list, false
		}
	}
	return
}

func getValueable(val interface{}) core.Valueable {
	if v, ok := val.(core.Valueable); ok {
		return v
	}
	return core.On(val)
}

// getValue returns the Value() of val iff val is a Valueable, else returns val
func getValue(val interface{}) interface{} {
	if v, ok := val.(core.Valueable); ok {
		return v.Value()
	}
	return val
}

func getLoop(v interface{}) (core.Valueable, bool) {
	if val, ok := v.(core.Valueable); ok {
		if _, ok := val.Value().(*core.Loop); ok {
			return val, true
		}
		return val, false
	}
	if l, ok := v.(*core.Loop); ok {
		return core.On(l), true
	}
	return core.On(v), false
}
