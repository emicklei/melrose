package dsl

import (
	"reflect"
	"testing"

	"github.com/emicklei/melrose/core"
)

func Test_envMap_Add(t *testing.T) {
	type args struct {
		l interface{}
		r interface{}
	}
	m := envMap{}
	s := NewVariableStore()
	s.Put("v1", 1)
	s.Put("v2", core.On(2))
	v1 := s.getVariable("v1")
	v2 := s.getVariable("v2")
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
		{
			"v2+1",
			m,
			args{v1, 1},
			2,
		},
		{
			"1+v2",
			m,
			args{1, v1},
			2,
		},
		{
			"v1+v2",
			m,
			args{v1, v2},
			3,
		},
		{
			"v2+v1",
			m,
			args{v2, v1},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Add(tt.args.l, tt.args.r); !reflect.DeepEqual(got.Value(), tt.want) {
				t.Errorf("envMap.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envMap_Sub(t *testing.T) {
	type args struct {
		l interface{}
		r interface{}
	}
	m := envMap{}
	s := NewVariableStore()
	s.Put("v1", 1)
	s.Put("v2", core.On(2))
	v1 := s.getVariable("v1")
	v2 := s.getVariable("v2")
	tests := []struct {
		name string
		e    envMap
		args args
		want interface{}
	}{
		{
			"v1-1",
			m,
			args{v1, 1},
			0,
		},
		{
			"1-v1",
			m,
			args{1, v1},
			0,
		},
		{
			"v1-v1",
			m,
			args{v1, v1},
			0,
		},
		{
			"v2-1",
			m,
			args{v2, 1},
			1,
		},
		{
			"1-v2",
			m,
			args{1, v2},
			-1,
		},
		{
			"v1-v2",
			m,
			args{v1, v2},
			-1,
		},
		{
			"v2-v1",
			m,
			args{v2, v1},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Sub(tt.args.l, tt.args.r); !reflect.DeepEqual(got.Value(), tt.want) {
				t.Errorf("envMap.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envMap_Multiply(t *testing.T) {
	type args struct {
		l interface{}
		r interface{}
	}
	m := envMap{}
	s := NewVariableStore()
	s.Put("v1", 1)
	s.Put("v2", core.On(2))
	v1 := s.getVariable("v1")
	v2 := s.getVariable("v2")
	tests := []struct {
		name string
		e    envMap
		args args
		want interface{}
	}{
		{
			"v1*1",
			m,
			args{v1, 1},
			1,
		},
		{
			"1*v1",
			m,
			args{1, v1},
			1,
		},
		{
			"v1*v1",
			m,
			args{v1, v1},
			1,
		},
		{
			"v2*1",
			m,
			args{v2, 1},
			2,
		},
		{
			"1*v2",
			m,
			args{1, v2},
			2,
		},
		{
			"v1*v2",
			m,
			args{v1, v2},
			2,
		},
		{
			"v2*v1",
			m,
			args{v2, v1},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Multiply(tt.args.l, tt.args.r); !reflect.DeepEqual(got.Value(), tt.want) {
				t.Errorf("envMap. Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}
