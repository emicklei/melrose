package calc

import (
	"reflect"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestAdd_Value(t *testing.T) {
	type fields struct {
		Left  interface{}
		Right interface{}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Add{
				Left:  tt.fields.Left,
				Right: tt.fields.Right,
			}
			if got := a.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
