package core

import (
	"reflect"
	"testing"
)

func Test_parseIndices(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			"just one",
			args{src: "1"},
			[][]int{{1}},
		},
		{
			"eleven",
			args{src: "11"},
			[][]int{{11}},
		},
		{
			"one [two three] four",
			args{src: "1 (2 3) 4"},
			[][]int{{1}, {2, 3}, {4}},
		},
		{
			"one two three",
			args{src: "1 2 3"},
			[][]int{{1}, {2}, {3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIndices(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIndices() = %v, want %v", got, tt.want)
			}
		})
	}
}
