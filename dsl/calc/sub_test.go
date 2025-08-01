package calc

import (
	"reflect"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestSub_Value(t *testing.T) {
	type fields struct {
		Left  any
		Right any
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{"1-2", fields{1, 2}, -1},
		{"1.0-2.0", fields{1.0, 2.0}, -1.0},
		{"1-[2]", fields{1, core.On(2)}, -1},
		{"[1]-[2]", fields{core.On(1), core.On(2)}, -1},
		{"[[1]]-[2]", fields{core.ValueHolder{Any: core.On(1)}, core.On(2)}, -1},
		{"1-nil", fields{1, nil}, 1},
		{"nil-2", fields{nil, 2}, -2},
		{"1.0-nil", fields{1.0, nil}, 1.0},
		{"nil-2.0", fields{nil, 2.0}, -2.0},
		{"nil-nil", fields{nil, nil}, 0},
		{"1.0-2", fields{1.0, 2}, -1.0},
		{"1-2.0", fields{1, 2.0}, -1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Sub{
				Left:  tt.fields.Left,
				Right: tt.fields.Right,
			}
			if got := a.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sub.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSub_Storex(t *testing.T) {
	s := Sub{Left: 1, Right: 2}
	if got, want := s.Storex(), "1 - 2"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
