package midi

import (
	"os"
	"testing"
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/m"
)

// BEAT=on go test -v -run "^(TestBeatmaster)$"
func TestBeatmaster(t *testing.T) {
	if on := os.Getenv("BEAT"); on != "on" {
		t.Skip()
	}
	dev, _ := Open()
	defer dev.Close()
	melrose.SetCurrentDevice(dev)
	dev.echo = false
	dev.bpm = 120

	b := melrose.NewBeatmaster(dev.bpm)
	b.Verbose(true)
	b.Start()
	time.Sleep(1 * time.Second)

	s1 := m.Sequence("8C 8E 8G")
	l1 := m.Loop(s1)
	b.Begin(l1)

	time.Sleep(3 * time.Second)

	s2 := m.Sequence("8F 8A 8C5")
	l2 := m.Loop(s2)
	b.Begin(l2)

	time.Sleep(10 * time.Second)
	b.End(l1)
	b.End(l2)

	time.Sleep(1 * time.Second)
	b.Stop()
}
