package calc

import (
	"reflect"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestMultiply_Value(t *testing.T) {
	type fields struct {
		Left  interface{}
		Right interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"3*2", fields{3, 2}, 6},
		{"3.0*2.0", fields{3.0, 2.0}, 6.0},
		{"3*[2]", fields{3, core.On(2)}, 6},
		{"[3]*[2]", fields{core.On(3), core.On(2)}, 6},
		{"[[3]]*[2]", fields{core.ValueHolder{Any: core.On(3)}, core.On(2)}, 6},
		{"3*nil", fields{3, nil}, 0},
		{"nil*2", fields{nil, 2}, 0},
		{"3.0*nil", fields{3.0, nil}, 0.0},
		{"nil*2.0", fields{nil, 2.0}, 0.0},
		{"nil*nil", fields{nil, nil}, 0},
		{"3.0*2", fields{3.0, 2}, 6.0},
		{"3*2.0", fields{3, 2.0}, 6.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Multiply{
				Left:  tt.fields.Left,
				Right: tt.fields.Right,
			}
			if got := a.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Multiply.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiply_Storex(t *testing.T) {
	m := Multiply{Left: 1, Right: 2}
	if got, want := m.Storex(), "1 * 2"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
