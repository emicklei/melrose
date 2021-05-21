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
