package core

import (
	"testing"
	"time"
)

func Test_tickerDuration(t *testing.T) {
	d := tickerDuration(60.0)
	if got, want := d, 1*time.Second; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	d2 := tickerDuration(300.0)
	if got, want := d2, 200*time.Millisecond; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
