package dsl

import (
	"reflect"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl/calc"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

type envMap map[string]interface{}

func exprOperators() []expr.Option {
	return []expr.Option{
		expr.Operator("-", "_Sub"),
		expr.Operator("+", "_Add"),
		expr.Operator("*", "_Multiply"),
		expr.Operator("<", "_LessThan"),
		expr.Operator("<=", "_LessEqualThan"),
		expr.Operator(">", "_GreaterThan"),
		expr.Operator(">=", "_GreaterEqualThan"),
		expr.Operator("!=", "_NotEqual"),
		expr.Operator("==", "_Equal"),
	}
}

func addOperatorsTo(env map[string]interface{}) {
	env["_Sub"] = envMap{}.Sub
	env["_Add"] = envMap{}.Add
	env["_Multiply"] = envMap{}.Multiply
	env["_LessThan"] = envMap{}.LessThan
	env["_LessEqualThan"] = envMap{}.LessEqualThan
	env["_GreaterThan"] = envMap{}.GreaterThan
	env["_GreaterEqualThan"] = envMap{}.GreaterEqualThan
	env["_NotEqual"] = envMap{}.NotEqual
	env["_Equal"] = envMap{}.Equal
}

func (envMap) Sub(l, r interface{}) core.HasValue {
	return calc.Sub{Left: l, Right: r}
}

func (envMap) Add(l, r any) core.HasValue {
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
	//log.Printf("%T %v\n", node, ast.Dump(*node))
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
		if n.Method {
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
		//log.Printf("%T %v %v\n", node, ast.Dump(*node), methodName)
	}
}
