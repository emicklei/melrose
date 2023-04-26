package dsl

import (
	"log"
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
		expr.Operator("*", "Multiply"),
		expr.Operator("<", "LessThan"),
		expr.Operator("<=", "LessEqualThan"),
		expr.Operator(">", "GreaterThan"),
		expr.Operator(">=", "GreaterEqualThan"),
		expr.Operator("!=", "NotEqual"),
		expr.Operator("==", "Equal"),
	}
}
func (envMap) Sub(l, r interface{}) core.HasValue {
	return calc.Sub{Left: l, Right: r}
}

func (envMap) Add(l, r interface{}) core.HasValue {
	return calc.Add{Left: l, Right: r}
}

func (envMap) Multiply(l, r interface{}) core.HasValue {
	return calc.Multiply{Left: l, Right: r}
}

func (envMap) LessThan(l, r interface{}) core.HasValue {
	return calc.NumberCompare{Left: l, Right: r, Operator: "<"}
}

func (envMap) LessEqualThan(l, r interface{}) core.HasValue {
	return calc.NumberCompare{Left: l, Right: r, Operator: "<="}
}

func (envMap) GreaterThan(l, r interface{}) core.HasValue {
	return calc.NumberCompare{Left: l, Right: r, Operator: ">"}
}

func (envMap) GreaterEqualThan(l, r interface{}) core.HasValue {
	return calc.NumberCompare{Left: l, Right: r, Operator: ">="}
}

func (envMap) NotEqual(l, r interface{}) core.HasValue {
	return calc.NumberCompare{Left: l, Right: r, Operator: "!="}
}

func (envMap) Equal(l, r interface{}) core.HasValue {
	return calc.NumberCompare{Left: l, Right: r, Operator: "=="}
}

var variableType = reflect.TypeOf(variable{})

// indexedAccessPatcher exist to patch expression which use [] on variables.
type indexedAccessPatcher struct{}

func (p *indexedAccessPatcher) Visit(node *ast.Node) {
	log.Printf("%T %v\n", node, ast.Dump(*node))
	n, ok := (*node).(*ast.MemberNode)
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
		_, ok = n.Property.(*ast.IdentifierNode)
		if ok {
			methodName = "AtVariable"
		}
		ast.Patch(node, &ast.CallNode{
			Callee: &ast.MemberNode{
				Node: in,
				Property: &ast.StringNode{
					Value: methodName,
				},
			},
			Arguments: []ast.Node{
				n.Property,
			},
		})
		log.Printf("%T %v %v\n", node, ast.Dump(*node), methodName)
	}
}
