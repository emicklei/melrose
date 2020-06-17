package file

import (
	"testing"
	"time"

	"github.com/emicklei/melrose"
)

func Test_microsecondsFromBPM(t *testing.T) {
	if got, want := quarterUSFromBPM(120.0), uint32(500000); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

}

func Test_ticksFromDuration(t *testing.T) {
	if got, want := ticksFromDuration(5*time.Second, quarterUSFromBPM(120.0)), uint32(9600); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func Test_Export(t *testing.T) {
	s := melrose.MustParseSequence("C")
	if err := Export("Test_Export.mid", s, 120.0); err != nil {
		t.Fatal(err)
	}
}
