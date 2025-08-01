package calc

import (
	"reflect"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestAdd_Value(t *testing.T) {
	type fields struct {
		Left  any
		Right any
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"1+2", fields{1, 2}, 3},
		{"1.0+2.0", fields{1.0, 2.0}, 3.0},
		{"1+[2]", fields{1, core.On(2)}, 3},
		{"[1]+[2]", fields{core.On(1), core.On(2)}, 3},
		{"[[1]]+[2]", fields{core.ValueHolder{Any: core.On(1)}, core.On(2)}, 3},
		{"1+nil", fields{1, nil}, 1},
		{"nil+2", fields{nil, 2}, 2},
		{"1.0+nil", fields{1.0, nil}, 1.0},
		{"nil+2.0", fields{nil, 2.0}, 2.0},
		{"nil+nil", fields{nil, nil}, 0},
		{"1.0+2", fields{1.0, 2}, 3.0},
		{"1+2.0", fields{1, 2.0}, 3.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Add{
				Left:  tt.fields.Left,
				Right: tt.fields.Right,
			}
			if got := a.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add.Value() = %v(%T), want %v(%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestAdd_Storex(t *testing.T) {
	a := Add{Left: 1, Right: 2}
	if got, want := a.Storex(), "1 + 2"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestAdd_S(t *testing.T) {
	n1, _ := core.ParseNote("C")
	n2, _ := core.ParseNote("D")
	a := Add{Left: n1, Right: n2}
	if got, want := a.S().Storex(), "sequence('C D')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	a = Add{Left: 1, Right: n2}
	if len(a.S().Notes) != 0 {
		t.Error("should be empty")
	}
	a = Add{Left: n1, Right: 1}
	if len(a.S().Notes) != 0 {
		t.Error("should be empty")
	}
}
