package core

import (
	"testing"
	"time"
)

func TestEuclideanRhythm(t *testing.T) {
	tests := []struct {
		name     string
		steps    int
		pulses   int
		rotation int
		expected []bool
	}{
		{
			name:     "basic 4 steps 1 pulse",
			steps:    4,
			pulses:   1,
			rotation: 0,
			expected: []bool{true, false, false, false},
		},
		{
			name:     "8 steps 3 pulses",
			steps:    8,
			pulses:   3,
			rotation: 0,
			expected: []bool{true, false, false, true, false, false, true, false},
		},
		{
			name:     "16 steps 5 pulses",
			steps:    16,
			pulses:   5,
			rotation: 0,
			expected: []bool{true, false, false, false, true, false, false, true, false, false, true, false, false, true, false, false},
		},
		{
			name:     "8 steps 5 pulses",
			steps:    8,
			pulses:   5,
			rotation: 0,
			expected: []bool{true, false, true, false, true, true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := euclideanRhythm(tt.steps, tt.pulses, tt.rotation)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("at index %d: expected %v, got %v", i, expected, result[i])
				}
			}
		})
	}
}

func TestRotatedRhythm(t *testing.T) {
	tests := []struct {
		name     string
		input    []bool
		rotation int
		expected []bool
	}{
		{
			name:     "no rotation",
			input:    []bool{true, false, true, false},
			rotation: 0,
			expected: []bool{true, false, true, false},
		},
		{
			name:     "rotate by 1",
			input:    []bool{true, false, true, false},
			rotation: 1,
			expected: []bool{false, true, false, true},
		},
		{
			name:     "rotate by 2",
			input:    []bool{true, false, true, false},
			rotation: 2,
			expected: []bool{true, false, true, false},
		},
		{
			name:     "rotate by 3",
			input:    []bool{true, false, true, false},
			rotation: 3,
			expected: []bool{false, true, false, true},
		},
		{
			name:     "negative rotation",
			input:    []bool{true, false, true, false},
			rotation: -1,
			expected: []bool{false, true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rotatedRhythm(tt.input, tt.rotation)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("at index %d: expected %v, got %v", i, expected, result[i])
				}
			}
		})
	}
}

func TestEuclidean_Storex(t *testing.T) {
	e := &Euclidean{
		Steps:    &ValueHolder{Any: 8},
		Beats:    &ValueHolder{Any: 3},
		Rotation: &ValueHolder{Any: 0},
		Playback: &ValueHolder{Any: "test"},
	}

	expected := "euclidean(8,3,0,'test')"
	if got := e.Storex(); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestEuclidean_Inspect(t *testing.T) {
	e := &Euclidean{
		Steps:    &ValueHolder{Any: 8},
		Beats:    &ValueHolder{Any: 3},
		Rotation: &ValueHolder{Any: 0},
		Playback: &ValueHolder{Any: "test"},
	}

	i := Inspection{Properties: make(map[string]any)}
	e.Inspect(i)

	// Check that all expected properties are set
	if i.Properties["steps"] != e.Steps {
		t.Errorf("expected steps to be %v, got %v", e.Steps, i.Properties["steps"])
	}
	if i.Properties["beats"] != e.Beats {
		t.Errorf("expected beats to be %v, got %v", e.Beats, i.Properties["beats"])
	}
	if i.Properties["rotation"] != e.Rotation {
		t.Errorf("expected rotation to be %v, got %v", e.Rotation, i.Properties["rotation"])
	}
	if i.Properties["playback"] != e.Playback {
		t.Errorf("expected playback to be %v, got %v", e.Playback, i.Properties["playback"])
	}

	// Check pattern generation
	pattern, ok := i.Properties["pattern"].(string)
	if !ok {
		t.Errorf("expected pattern to be a string, got %T", i.Properties["pattern"])
	}
	expectedPattern := "!..!..!." // 8 steps, 3 beats
	if pattern != expectedPattern {
		t.Errorf("expected pattern %s, got %s", expectedPattern, pattern)
	}
}

func TestEuclidean_Handle(t *testing.T) {
	e := &Euclidean{}
	timeline := &Timeline{}
	when := time.Now()

	// Handle should not panic and should be a no-op
	e.Handle(timeline, when)
}

func TestEuclideanRhythm_WithRotation(t *testing.T) {
	tests := []struct {
		name     string
		steps    int
		pulses   int
		rotation int
		expected []bool
	}{
		{
			name:     "4 steps 2 pulses rotation 1",
			steps:    4,
			pulses:   2,
			rotation: 1,
			expected: []bool{false, true, false, true},
		},
		{
			name:     "8 steps 3 pulses rotation 2",
			steps:    8,
			pulses:   3,
			rotation: 2,
			expected: []bool{true, false, true, false, false, true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := euclideanRhythm(tt.steps, tt.pulses, tt.rotation)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("at index %d: expected %v, got %v", i, expected, result[i])
				}
			}
		})
	}
}

func TestEuclidean_InspectPatternGeneration(t *testing.T) {
	tests := []struct {
		name            string
		steps           int
		beats           int
		rotation        int
		expectedPattern string
	}{
		{
			name:            "4 steps 1 beat",
			steps:           4,
			beats:           1,
			rotation:        0,
			expectedPattern: "!...",
		},
		{
			name:            "8 steps 3 beats",
			steps:           8,
			beats:           3,
			rotation:        0,
			expectedPattern: "!..!..!.",
		},
		{
			name:            "8 steps 5 beats",
			steps:           8,
			beats:           5,
			rotation:        0,
			expectedPattern: "!.!.!!.!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Euclidean{
				Steps:    &ValueHolder{Any: tt.steps},
				Beats:    &ValueHolder{Any: tt.beats},
				Rotation: &ValueHolder{Any: tt.rotation},
				Playback: &ValueHolder{Any: "test"},
			}

			i := Inspection{Properties: make(map[string]any)}
			e.Inspect(i)

			pattern, ok := i.Properties["pattern"].(string)
			if !ok {
				t.Errorf("expected pattern to be a string, got %T", i.Properties["pattern"])
			}
			if pattern != tt.expectedPattern {
				t.Errorf("expected pattern %s, got %s", tt.expectedPattern, pattern)
			}
		})
	}
}
