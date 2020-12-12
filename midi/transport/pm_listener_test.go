// +build !udp

package transport

import (
	"testing"

	"github.com/emicklei/melrose/control"
)

func Test_listener_remove(t *testing.T) {
	lis := newListener(nil)
	usr := new(control.Listen)
	lis.Add(usr)
	lis.Remove(usr)
	if got, want := len(lis.noteListeners), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
