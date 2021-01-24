package dsl

import (
	"reflect"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl/calc"
)

type envMap map[string]interface{}

func (envMap) exprOperators() []expr.Option {
	return []expr.Option{
		expr.Operator("-", "Sub"),
		expr.Operator("+", "Add"),
		expr.Operator("*", "Mulitply"),
	}
}
func (envMap) Sub(l, r interface{}) core.Valueable {
	return calc.Sub{Left: l, Right: r}
	// if vl, ok := l.(variable); ok {
	// 	return vl.dispatchSub(r)
	// }
	// if vr, ok := r.(variable); ok {
	// 	return vr.dispatchSubFrom(l)
	// }
	// if li, ok := l.(int); ok {
	// 	if ri, ok := r.(int); ok {
	// 		return li - ri
	// 	}
	// }
	// notify.Panic(fmt.Errorf("substraction failed [%v (%T) - %v (%T)]", l, l, r, r))
	// return nil
}

func (envMap) Add(l, r interface{}) core.Valueable {
	return calc.Add{Left: l, Right: r}
	// if vl, ok := l.(variable); ok {
	// 	return vl.dispatchAdd(r)
	// }
	// if vr, ok := r.(variable); ok {
	// 	return vr.dispatchAdd(l)
	// }
	// if li, ok := l.(int); ok {
	// 	if ri, ok := r.(int); ok {
	// 		return li + ri
	// 	}
	// }
	// notify.Panic(fmt.Errorf("addition failed [%v (%T) + %v (%T)]", l, l, r, r))
	//return nil
}

func (envMap) Mulitply(l, r interface{}) core.Valueable {
	return calc.Multiply{Left: l, Right: r}
	// if vl, ok := l.(variable); ok {
	// 	return vl.dispatchMultiply(r)
	// }
	// if vr, ok := r.(variable); ok {
	// 	return vr.dispatchMultiply(l)
	// }
	// if li, ok := l.(int); ok {
	// 	if ri, ok := r.(int); ok {
	// 		return li * ri
	// 	}
	// }
	// notify.Panic(fmt.Errorf("multiplication failed [%v (%T) + %v (%T)]", l, l, r, r))
	// return nil
}

var variableType = reflect.TypeOf(variable{})

// indexedAccessPatcher exist to patch expression which use [] on variables.
type indexedAccessPatcher struct{}

func (p *indexedAccessPatcher) Enter(_ *ast.Node) {}
func (p *indexedAccessPatcher) Exit(node *ast.Node) {
	n, ok := (*node).(*ast.IndexNode)
	if ok {
		// check receiver type
		in, ok := n.Node.(*ast.IdentifierNode)
		if !ok {
			return
		}
		if in.Type() != variableType {
			return
		}
		// check argument type
		methodName := "At"
		in, ok = n.Index.(*ast.IdentifierNode)
		if ok {
			methodName = "AtVariable"
		}
		ast.Patch(node, &ast.MethodNode{
			Node:   n.Node,
			Method: methodName,
			Arguments: []ast.Node{
				n.Index,
			},
		})
	}
}
