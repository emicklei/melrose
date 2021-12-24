package core

import (
	"testing"
	"time"
)

/**
    play_test.go:13: bpm 120
    play_test.go:15: whole 2s
    play_test.go:19: 1C 2s
    play_test.go:19: 2C 1s
    play_test.go:19: C 500ms
    play_test.go:19: 8C 250ms
	play_test.go:19: 16C 125ms
**/
func TestDurationToFraction(t *testing.T) {
	type args struct {
		bpm float64
		d   time.Duration
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{"250ms", args{120.0, 250 * time.Millisecond}, 0.125},
		{"100ms", args{120.0, 100 * time.Millisecond}, 0.0625},
		{"175ms", args{120.0, 175 * time.Millisecond}, 0.0625},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationToFraction(tt.args.bpm, tt.args.d); got != tt.want {
				t.Errorf("DurationToFraction() = %v, want %v", got, tt.want)
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
