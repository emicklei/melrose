package calc

import (
	"reflect"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestCompare_Value(t *testing.T) {
	type fields struct {
		Left     any
		Right    any
		Operator string
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{"1<2", fields{1, 2, "<"}, true},
		{"1.0>2.0", fields{1.0, 2.0, ">"}, false},
		{"3.0==3.0", fields{3.0, 3.0, "=="}, true},
		{"3.0==3", fields{3.0, 3, "=="}, true},
		{"1<=[2]", fields{1, core.On(2), "<="}, true},
		{"[1]>=[2]", fields{core.On(1), core.On(2), ">="}, false},
		{"[[1]]!=[2]", fields{core.ValueHolder{Any: core.On(1)}, core.On(2), "!="}, true},
		{"1<nil", fields{1, nil, "<"}, false},
		{"nil<1", fields{nil, 1, "<"}, true},
		{"1.0<nil", fields{1.0, nil, "<"}, false},
		{"nil<1.0", fields{nil, 1.0, "<"}, true},
		{"1==1", fields{1, 1, "=="}, true},
		{"1>=1", fields{1, 1, ">="}, true},
		{"1<=1", fields{1, 1, "<="}, true},
		{"1!=2", fields{1, 2, "!="}, true},
		{"1>0", fields{1, 0, ">"}, true},
		{"1<0", fields{1, 0, "<"}, false},
		{"1==bad", fields{1, 0, "bad"}, false},
		{"1.0==bad", fields{1.0, 0, "bad"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NumberCompare{
				Left:     tt.fields.Left,
				Right:    tt.fields.Right,
				Operator: tt.fields.Operator,
			}
			if got := a.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberCompare.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberCompare_Storex(t *testing.T) {
	nc := NumberCompare{Left: 1, Right: 2, Operator: "<"}
	if got, want := nc.Storex(), "1 < 2"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
