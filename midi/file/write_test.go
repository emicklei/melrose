package file

import (
	"github.com/emicklei/melrose/core"
	"testing"
	"time"
)

func Test_microsecondsFromBPM(t *testing.T) {
	if got, want := quarterUSFromBPM(120.0), uint32(500000); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

}

func Test_ticksFromDuration(t *testing.T) {
	for _, each := range []struct {
		dur   string
		bpm   float64
		ticks int
	}{
		{
			"5s",
			120.0,
			9600,
		},
		{
			"1s",
			120.0,
			1920,
		},
		{
			"125ms",
			120.0,
			240,
		},
	} {
		d, _ := time.ParseDuration(each.dur)
		if got, want := ticksFromDuration(d, quarterUSFromBPM(each.bpm)), uint32(each.ticks); got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func Test_Export(t *testing.T) {
	s := core.MustParseSequence("C")
	if err := Export("Test_Export.mid", s, 120.0); err != nil {
		t.Fatal(err)
	}
}
