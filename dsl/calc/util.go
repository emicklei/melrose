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

func resolveFloat(v interface{}) (float64, bool) {
	if i, ok := v.(float64); ok {
		return i, true
	}
	if v, ok := v.(core.HasValue); ok {
		return resolveFloat(v.Value())
	}
	return 0.0, false
}
