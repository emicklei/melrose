package core

import (
	"testing"
	"time"
)

func TestMIDI_Failures(t *testing.T) {
	{
		m := MIDINote{duration: On(-1), number: On(60), velocity: On(60)}
		if _, err := m.ToNote(); err == nil {
			t.Fail()
		}
	}
	{
		m := MIDINote{duration: On(500), number: On(-1), velocity: On(60)}
		if _, err := m.ToNote(); err == nil {
			t.Fail()
		}
	}
	{
		m := MIDINote{duration: On(500), number: On(60), velocity: On(-1)}
		if _, err := m.ToNote(); err == nil {
			t.Fail()
		}
	}
}

func TestMIDI_ToNote(t *testing.T) {
	type fields struct {
		duration Valueable
		number   Valueable
		velocity Valueable
	}
	tests := []struct {
		name     string
		fields   fields
		want     Note
		duration time.Duration
	}{
		{
			"F+",
			fields{On(4), On(65), On(64)},
			MustParseNote("F+"),
			ZeroDuration,
		},
		{
			"G3+",
			fields{On(8), On(55), On(64)},
			MustParseNote("8G3+"),
			ZeroDuration,
		},
		{
			"E♭5",
			fields{On(16), On(75), On(Normal)},
			MustParseNote("16E♭5"),
			ZeroDuration,
		},
		{
			"C-int",
			fields{On(600), On(60), On(Normal)},
			MustParseNote("C"),
			time.Duration(600) * time.Millisecond,
		},
		{
			"C-time",
			fields{On(time.Duration(600) * time.Millisecond), On(60), On(Normal)},
			MustParseNote("C"),
			time.Duration(600) * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MIDINote{
				duration: tt.fields.duration,
				number:   tt.fields.number,
				velocity: tt.fields.velocity,
			}
			n, _ := m.ToNote()
			// if got, want := n.duration, time.Duration(500)*time.Millisecond; got != want {
			// 	t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
			// }
			if got, want := n.Name, tt.want.Name; got != want {
				t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
			}
			if got, want := n.Velocity, tt.want.Velocity; got != want {
				t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
			}
			if got, want := n.DurationFactor(), tt.want.DurationFactor(); got != want {
				t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
			}
			if got, want := n.duration, tt.duration; got != want {
				t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
			}
		})
	}
}
