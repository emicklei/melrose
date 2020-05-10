package dsl

import (
	"reflect"
	"testing"
)

func Test_envMap_Add(t *testing.T) {
	type args struct {
		l interface{}
		r interface{}
	}
	m := envMap{}
	s := NewVariableStore()
	s.Put("v1", 1)
	v1 := variable{Name: "v1", store: s}
	tests := []struct {
		name string
		e    envMap
		args args
		want interface{}
	}{
		{
			"v1+1",
			m,
			args{v1, 1},
			2,
		},
		{
			"1+v1",
			m,
			args{1, v1},
			2,
		},
		{
			"v1+v1",
			m,
			args{v1, v1},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Add(tt.args.l, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("envMap.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
