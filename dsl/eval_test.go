package dsl

import (
	"github.com/emicklei/melrose/core"
	"testing"

	"github.com/emicklei/melrose/op"
)

func TestIsCompatible(t *testing.T) {
	if got, want := true, IsCompatibleSyntax("1.0"); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := false, IsCompatibleSyntax("2.0"); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := false, IsCompatibleSyntax("1.1"); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNestedFunctions(t *testing.T) {
	e := NewEvaluator(testContext())
	input := `pitch(1,repeat(1,reverse(join(note('E'),sequence('F G')))))`
	_, err := e.EvaluateExpression(input)
	if err != nil {
		t.Error(err)
	}
}

func TestMulitLineEvaluate(t *testing.T) {
	e := NewEvaluator(testContext())
	input := `sequence("
		C D E C 
		C D E C 
		E F 2G
		E F 2G 
		8G 8A 8G 8F E C 
		8G 8A 8G 8F E C
		2C 2G3 2C
		2C 2G3 2C
		")`
	_, err := e.EvaluateStatement(input)
	if err != nil {
		t.Error(err)
	}
}

func Test_isAssignment(t *testing.T) {
	type args struct {
		entry string
	}
	tests := []struct {
		name           string
		args           args
		wantVarname    string
		wantExpression string
		wantOk         bool
	}{
		{"a=1",
			args{"a=1"},
			"a",
			"1",
			true,
		},
		{" a = note('=')",
			args{" a = note('=')"},
			"a",
			"note('=')",
			true,
		},
		{"multi line",
			args{`j2 = join(  repeat(2,i2) )`},
			"j2",
			"join(  repeat(2,i2) )",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVarname, gotExpression, gotOk := IsAssignment(tt.args.entry)
			if gotVarname != tt.wantVarname {
				t.Errorf("isAssignment() gotVarname = %v, want %v", gotVarname, tt.wantVarname)
			}
			if gotExpression != tt.wantExpression {
				t.Errorf("isAssignment() gotExpression = %v, want %v", gotExpression, tt.wantExpression)
			}
			if gotOk != tt.wantOk {
				t.Errorf("isAssignment() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestEvaluateProgram_SingleLine(t *testing.T) {
	e := NewEvaluator(testContext())
	if _, err := e.EvaluateProgram(`a = 1`); err != nil {
		t.Error(err)
	}
}

func TestEvaluateProgram_CommentLine(t *testing.T) {
	e := NewEvaluator(testContext())
	if _, err := e.EvaluateProgram(`// a = 1`); err != nil {
		t.Error(err)
	}
}

func TestEvaluateProgram_FirstTab(t *testing.T) {
	e := NewEvaluator(testContext())
	if _, err := e.EvaluateProgram(`	a = 1`); err == nil {
		t.Error(err)
	}
}

func testContext() core.Context {
	return core.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     core.NoLooper,
	}
}

func TestEvaluateProgram_TwoLines(t *testing.T) {
	e := NewEvaluator(testContext())
	r, err := e.EvaluateProgram(
		`a = 1
b = 2`)
	if err != nil {
		t.Error(err)
	}
	if got, want := r, 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEvaluateProgram_BrokenSequence(t *testing.T) {
	e := NewEvaluator(testContext())
	r, err := e.EvaluateProgram(
		`a = sequence(
	'A')`)
	if err != nil {
		t.Error(err)
	}
	if got, want := r, 1; got.(core.Sequence).Length() != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestEvaluateProgram_AllBrokenSequence(t *testing.T) {
	e := NewEvaluator(testContext())
	r, err := e.EvaluateProgram(
		`a = sequence
	(
	'
	A
	'
	)`)
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Fatal()
	}
	if got, want := r, 1; got.(core.Sequence).Length() != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEvaluateProgram_TwoTabs(t *testing.T) {
	e := NewEvaluator(testContext())
	r, err := e.EvaluateProgram(
		`a
	=
		1`)
	if err != nil {
		t.Error(err)
	}
	if got, want := r, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEvaluateProgram_TrailingWhitespace(t *testing.T) {
	e := newTestEvaluator()
	r, err := e.EvaluateProgram(
		`a=1
 `)
	checkError(t, err)
	if got, want := r, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEvaluateError_Play(t *testing.T) {
	_, err := newTestEvaluator().evaluateCleanStatement(`repeat(-1,1)`)
	if err == nil {
		t.Fail()
	}
}

func TestEvaluateIndexOnVariable(t *testing.T) {
	e := newTestEvaluator()
	r, err := e.EvaluateProgram(
		`([1,2,3])[1]
a=[1,2,3]
a[1]`)
	checkError(t, err)
	if got, want := r, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEvaluate_Scale_At(t *testing.T) {
	e := newTestEvaluator()
	r, err := e.EvaluateProgram(
		`at(1,scale(2,'C'))`)
	checkError(t, err)
	at, _ := r.(op.AtIndex)
	if got, want := at.Target.S().At(0)[0], core.MustParseNote("C"); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
