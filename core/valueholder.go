package core

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/notify"
)

func String(h Valueable) string {
	if h == nil {
		return ""
	}
	val := h.Value()
	if val == nil {
		return ""
	}
	if v, ok := val.(string); ok {
		return v
	}
	// maybe the value is a Valueable
	if vv, ok := val.(Valueable); ok {
		return String(vv)
	}
	return ""
}

func Float(h Valueable) float32 {
	if h == nil {
		return 0.0
	}
	val := h.Value()
	if val == nil {
		return 0.0
	}
	if v, ok := val.(float32); ok {
		return v
	}
	if v, ok := val.(float64); ok {
		return float32(v)
	}
	if v, ok := val.(int); ok {
		return float32(v)
	}
	// maybe the value is a Valueable
	if vv, ok := val.(Valueable); ok {
		return Float(vv)
	}
	return 0.0
}

func Duration(h Valueable) time.Duration {
	if h == nil {
		return time.Duration(0)
	}
	val := h.Value()
	if val == nil {
		return time.Duration(0)
	}
	if v, ok := val.(int); ok {
		return time.Duration(v) * time.Millisecond
	}
	if v, ok := val.(time.Duration); ok {
		return v
	}
	// maybe the value is a Valueable
	if vv, ok := val.(Valueable); ok {
		return Duration(vv)
	}
	notify.Print(notify.Warningf("Duration() expected [time.Duration|int] but got [%T], return 0", h.Value()))
	return time.Duration(0)
}

func Int(h Valueable) int {
	return getInt(h, false)
}

func getInt(h Valueable, silent bool) int {
	if h == nil {
		return 0
	}
	val := h.Value()
	if val == nil {
		return 0
	}
	if v, ok := val.(int); ok {
		return v
	}
	// maybe the value is a Valueable
	if vv, ok := val.(Valueable); ok {
		return getInt(vv, silent)
	}
	if !silent {
		notify.Print(notify.Warningf("Int() expected [int] but got [%T], return 0", h.Value()))
	}
	return 0
}

func ToValueable(v interface{}) Valueable {
	if w, ok := v.(Valueable); ok {
		return w
	}
	return &ValueHolder{Any: v}
}

func On(v interface{}) Valueable {
	return ToValueable(v)
}

// ValueHolder is decorate any object to become a Valueable.
type ValueHolder struct {
	Any interface{}
}

func (h ValueHolder) Value() interface{} {
	return h.Any
}

func (h ValueHolder) Storex() string {
	if st, ok := h.Any.(Storable); ok {
		return st.Storex()
	}
	if s, ok := h.Any.(string); ok {
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("%v", h.Any)
}

func (h ValueHolder) String() string {
	return h.Storex()
}
