package core

import (
	"testing"
	"time"
)

/*
*

	    play_test.go:13: bpm 120
	    play_test.go:15: whole 2s
	    play_test.go:19: 1C 2s
	    play_test.go:19: 2C 1s
	    play_test.go:19: C 500ms
	    play_test.go:19: 8C 250ms
		play_test.go:19: 16C 125ms

*
*/
func TestDurationToFraction(t *testing.T) {
	type args struct {
		bpm float64
		d   time.Duration
	}
	tests := []struct {
		name       string
		args       args
		want       float32
		wantDotted bool
	}{
		{"2s", args{120.0, 2 * time.Second}, 1.0, false},
		{"250ms", args{120.0, 250 * time.Millisecond}, 0.125, false},
		{"375ms", args{120.0, 375 * time.Millisecond}, 0.125, true},
		{"100ms", args{120.0, 100 * time.Millisecond}, 0.03125, true},
		{"175ms", args{120.0, 175 * time.Millisecond}, 0.0625, true},
		{"60ms", args{120.0, 60 * time.Millisecond}, 0.03125, false},
		{"quarter at 60bpm", args{60.0, time.Second}, 0.25, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, dotted := DurationToFraction(tt.args.bpm, tt.args.d)
			if got != tt.want || dotted != tt.wantDotted {
				t.Errorf("DurationToFraction() = %v,%v want %v,%v", got, dotted, tt.want, tt.wantDotted)
			}
		})
	}
}
