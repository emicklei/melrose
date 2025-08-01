package calc

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func Test_resolveFloatWithInt(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 bool
	}{
		{"float", args{1.0}, 1.0, true},
		{"int", args{1}, 1.0, true},
		{"nil", args{nil}, 0.0, true},
		{"hasvalue", args{core.On(1.0)}, 1.0, true},
		{"string", args{"string"}, 0.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := resolveFloatWithInt(tt.args.v)
			if got != tt.want {
				t.Errorf("resolveFloatWithInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("resolveFloatWithInt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_resolveInt(t *testing.T) {
	if _, ok := resolveInt(1.1); ok {
		t.Fail()
	}
}
