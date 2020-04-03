package dsl

import (
	"testing"
)

func TestNestedFunctions(t *testing.T) {
	input := `pitch(1,repeat(1,reverse(join(note('E'),sequence('F G')))))`
	v, err := Evaluate(NewVariableStore(), input)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", v)
}

func TestMulitLineEvaluate(t *testing.T) {
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
	t.Log(input)
	_, err := Evaluate(NewVariableStore(), input)
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
