package melrose

import (
	"fmt"

	"github.com/emicklei/melrose/notify"
)

type Valueable interface {
	Value() interface{}
}

func Int(h Valueable) int {
	if v, ok := h.Value().(int); ok {
		return v
	}
	// maybe the value is a Valueable
	if vv, ok := h.Value().(Valueable); ok {
		return Int(vv)
	}
	notify.Print(notify.Warningf("expected [int] but got [%T]", h.Value()))
	return 0
}

func ToValueable(v interface{}) Valueable {
	if w, ok := v.(Valueable); ok {
		return w
	}
	return &ValueHolder{Any: v}
}

// ValueHolder is decorate any object to become a Valueable.
type ValueHolder struct {
	Any interface{}
}

func (h ValueHolder) Value() interface{} {
	return h.Any
}

func (h *ValueHolder) SetValue(newAny interface{}) {
	h.Any = newAny
}

func (h ValueHolder) Storex() string {
	return fmt.Sprintf("%v", h.Any)
}

func (h ValueHolder) String() string {
	return h.Storex()
}

func On(v interface{}) *ValueHolder {
	return &ValueHolder{Any: v}
}
