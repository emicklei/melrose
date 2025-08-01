package calc

import "github.com/emicklei/melrose/core"

func resolveInt(v interface{}) (int, bool) {
	if i, ok := v.(int); ok {
		return i, true
	}
	if v, ok := v.(core.HasValue); ok {
		return resolveInt(v.Value())
	}
	return 0, false
}

func resolveFloatWithInt(v interface{}) (float64, bool) {
	if f, ok := v.(float64); ok {
		return f, true
	}
	if i, ok := v.(int); ok {
		return float64(i), true
	}
	if v == nil {
		return 0.0, true
	}
	if v, ok := v.(core.HasValue); ok {
		return resolveFloatWithInt(v.Value())
	}
	return 0.0, false
}
