package dsl

import (
	"reflect"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
)

type envMap map[string]interface{}

func (envMap) exprOperators() []expr.Option {
	return []expr.Option{
		expr.Operator("-", "Sub"),
		expr.Operator("+", "Add"),
	}
}
func (envMap) Sub(v variable, i int) int { return i }

func (envMap) Add(l, r interface{}) interface{} {
	if vl, ok := l.(variable); ok {
		return vl.dispatchAdd(r)
	}
	if vr, ok := r.(variable); ok {
		return vr.dispatchAdd(l)
	}
	return nil
}

func (v variable) dispatchAdd(r interface{}) interface{} {
	if vr, ok := r.(variable); ok {
		// int
		il, lok := v.Value().(int)
		ir, rok := vr.Value().(int)
		if lok && rok {
			return il + ir
		}
	}
	if ir, ok := r.(int); ok {
		il, lok := v.Value().(int)
		if lok {
			return il + ir
		}
	}
	return nil
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
