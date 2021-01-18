package op

import (
	"reflect"
	"testing"
)

func Test_parseIndexOffsets(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		wantM []int2int
	}{
		{"one to one", args{"1:1"}, []int2int{{1, 1}}},
		{"comma separated", args{"1:1,2:2"}, []int2int{{1, 1}, {2, 2}}},
		{"extra spaces", args{" -1:1 , 2:-2, 3:3, 4:4"}, []int2int{{-1, 1}, {2, -2}, {3, 3}, {4, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotM := parseIndexOffsets(tt.args.s); !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("parseIndexOffsets() = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}

func Test_parseIndexFloats(t *testing.T) {
	m := parseIndexFloats("1:1, 2:1.0, 3:0.5, 4:0.01625, 1:2, 1:4, 1:8, 1:16")
	if got, want := m[0].at, 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[0].float, float32(1.0); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[1].float, float32(1.0); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[2].float, float32(0.5); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[3].float, float32(0.01625); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
