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
		{"3*[2]", fields{3, core.On(2)}, 6},
		{"[3]*[2]", fields{core.On(3), core.On(2)}, 6},
		{"[[3]]*[2]", fields{core.ValueHolder{Any: core.On(3)}, core.On(2)}, 6},
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
