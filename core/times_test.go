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

func TestFractionToDurationParts(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name         string
		args         args
		wantFraction float32
		wantDotted   bool
	}{
		{"1.0", args{1.0}, 1.0, false},
		{"0.8", args{0.8}, 0.5, true},
		{"0.7", args{0.7}, 0.5, true},
		{"0.6", args{0.6}, 0.5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFraction, gotDotted := FractionToDurationParts(tt.args.f)
			if gotFraction != tt.wantFraction {
				t.Errorf("FractionToDurationParts(%s) gotFraction = %v, want %v", tt.name, gotFraction, tt.wantFraction)
			}
			if gotDotted != tt.wantDotted {
				t.Errorf("FractionToDurationParts(%s) gotDotted = %v, want %v", tt.name, gotDotted, tt.wantDotted)
			}
		})
	}
}
