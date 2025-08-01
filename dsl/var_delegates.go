package dsl

import (
	"fmt"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
)

// At is called from expr after patching []. One-based
func (v variable) At(index int) any {
	m, ok := v.store.Get(v.Name)
	if !ok {
		return nil
	}
	if intArray, ok := m.([]any); ok {
		if index < 1 || index > len(intArray) {
			return nil
		}
		return intArray[index-1]
	}
	if indexable, ok := m.(core.Indexable); ok {
		return indexable.At(index)
	}
	if sequenceable, ok := m.(core.Sequenceable); ok {
		return core.BuildSequence(sequenceable.S().At(index))
	}
	return nil
}

// AtVariable is called from expr after patching [].
func (v variable) AtVariable(index variable) any {
	indexVal := core.Int(index)
	if indexVal == 0 {
		return nil
	}
	return v.At(indexVal)
}

// dispatchSubFrom  v(l) - r
func (v variable) dispatchSub(r any) any {
	if vr, ok := r.(core.HasValue); ok {
		// int
		il, lok := resolveInt(v)
		ir, rok := resolveInt(vr)
		if lok && rok {
			return il - ir
		}
	}
	if ir, ok := r.(int); ok {
		// int
		il, lok := resolveInt(v)
		if lok {
			return il - ir
		}
	}
	notify.Panic(fmt.Errorf("subtraction failed [%v (%T) - %v (%T)]", v, v, r, r))
	return nil
}

// dispatchSubFrom  l - v(r)
func (v variable) dispatchSubFrom(l any) any {
	if vl, ok := l.(core.HasValue); ok {
		// int
		il, lok := resolveInt(vl)
		ir, rok := resolveInt(v)
		if lok && rok {
			return il - ir
		}
	}
	if il, ok := l.(int); ok {
		// int
		ir, rok := resolveInt(v)
		if rok {
			return il - ir
		}
	}
	notify.Panic(fmt.Errorf("subtraction failed [%v (%T) - %v (%T)]", l, l, v, v))
	return nil
}

func (v variable) dispatchAdd(r any) any {
	if vr, ok := r.(core.HasValue); ok {
		// int
		il, lok := resolveInt(v)
		ir, rok := resolveInt(vr)
		if lok && rok {
			return il + ir
		}
	}
	if ir, ok := r.(int); ok {
		il, lok := resolveInt(v)
		if lok {
			return il + ir
		}
	}
	notify.Panic(fmt.Errorf("addition failed [%v (%T) + %v (%T)]", r, r, v, v))
	return nil
}

// func (v variable) dispatchMultiply(r any) any {
// 	if vr, ok := r.(core.HasValue); ok {
// 		// int
// 		il, lok := resolveInt(v)
// 		ir, rok := resolveInt(vr)
// 		if lok && rok {
// 			return il * ir
// 		}
// 	}
// 	if ir, ok := r.(int); ok {
// 		il, lok := resolveInt(v)
// 		if lok {
// 			return il * ir
// 		}
// 	}
// 	notify.Panic(fmt.Errorf("multiplication failed [%v (%T) * %v (%T)]", r, r, v, v))
// 	return nil
// }

func resolveInt(v core.HasValue) (int, bool) {
	vv := v.Value()
	if i, ok := vv.(int); ok {
		return i, true
	}
	if vvv, ok := vv.(core.HasValue); ok {
		return resolveInt(vvv)
	}
	return 0, false
}
