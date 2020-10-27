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
