package transport

import (
	"errors"
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func UseUDPTransport(port int) {
	if core.IsDebug() {
		notify.Debugf("transport.UseUDPTransport with [:%d]", port)
	}
	Factory = func() Transporter {
		return &UDPTransporter{port: port}
	}
}

type UDPTransporter struct {
	port int
}

func (t *UDPTransporter) NewMIDIOut(id int) (MIDIOut, error) {
	return newRouterClient(t.port, id)
}

func (t *UDPTransporter) NewMIDIIn(id int) (MIDIIn, error) {
	return nil, errors.New("input unsupported")
}

func (t *UDPTransporter) HasInputCapability() bool {
	return false
}

func (t *UDPTransporter) Terminate() {

}

func (t *UDPTransporter) DefaultOutputDeviceID() int {
	return 1
}

func (t *UDPTransporter) NewMIDIListener(in MIDIIn) MIDIListener {
	return nil
}

func (t *UDPTransporter) PrintInfo(inID, outID int) {
	notify.PrintHighlighted("Usage:")
	fmt.Println()
}
