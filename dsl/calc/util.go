package calc

import "github.com/emicklei/melrose/core"

func resolveInt(v interface{}) (int, bool) {
	if i, ok := v.(int); ok {
		return i, true
	}
	if v, ok := v.(core.Valueable); ok {
		return resolveInt(v.Value())
	}
	return 0, false
}
