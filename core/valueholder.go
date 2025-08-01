package core

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/notify"
)

func String(h HasValue) string {
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
	// maybe the value is a HasValue
	if vv, ok := val.(HasValue); ok {
		return String(vv)
	}
	return ""
}

func Float(h HasValue) float32 {
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
	// maybe the value is a HasValue
	if vv, ok := val.(HasValue); ok {
		return Float(vv)
	}
	return 0.0
}

func Duration(h HasValue) time.Duration {
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
	// maybe the value is a HasValue
	if vv, ok := val.(HasValue); ok {
		return Duration(vv)
	}
	notify.Warnf("Duration() expected [time.Duration|int] but got [%T], return 0", h.Value())
	return time.Duration(0)
}

func Int(h HasValue) int {
	return getInt(h, false)
}

func ToSequenceable(v HasValue) Sequenceable {
	if v == nil {
		return EmptySequence
	}
	val := v.Value()
	if val == nil {
		return EmptySequence
	}
	if seq, ok := val.(Sequenceable); ok {
		return seq
	}
	return EmptySequence
}

func getInt(h HasValue, silent bool) int {
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
	// maybe the value is a HasValue
	if vv, ok := val.(HasValue); ok {
		return getInt(vv, silent)
	}
	if !silent {
		notify.Warnf("Int() expected [int] but got [%T], return 0", h.Value())
	}
	return 0
}

func ToHasValue(v any) HasValue {
	if w, ok := v.(HasValue); ok {
		return w
	}
	return &ValueHolder{Any: v}
}

// ValueOf returns the non HasValue value of v
func ValueOf(v any) any {
	if w, ok := v.(HasValue); ok {
		return ValueOf(w.Value())
	}
	return v
}

// IndexOf returns the non HasValue value of v
func IndexOf(v any) any {
	if i, ok := v.(HasIndex); ok {
		return i.Index()
	}
	if w, ok := v.(HasValue); ok {
		return IndexOf(w.Value())
	}
	return 0 // no index
}

func On(v any) HasValue {
	return ToHasValue(v)
}

// ValueHolder is decorate any object to become a HasValue.
type ValueHolder struct {
	Any any
}

func (h ValueHolder) Value() any {
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

// TODO used?
func (h ValueHolder) String() string {
	return h.Storex()
}

type ValueFunction struct {
	StoreString string
	Function    func() any
}

func (v ValueFunction) Storex() string {
	return v.StoreString
}

func (v ValueFunction) Value() any {
	return v.Function()
}
