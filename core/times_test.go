package core

import (
	"testing"
	"time"
)

func TestWholeNoteDuration(t *testing.T) {
	type args struct {
		bpm float64
	}
	s2, _ := time.ParseDuration("2s")
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "120",
			args: args{bpm: 120.0},
			want: s2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WholeNoteDuration(tt.args.bpm); got != tt.want {
				t.Errorf("WholeNoteDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
