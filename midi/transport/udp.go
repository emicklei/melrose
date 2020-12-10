package transport

import (
	"errors"
	"fmt"
)

func UseUDPTransport(hostport string) {
	Factory = func() Transporter {
		return &UDPTransporter{hostport: hostport}
	}
}

type UDPTransporter struct {
	hostport string
}

func (t *UDPTransporter) NewMIDIOut(id int) (MIDIOut, error) {
	return newRouterClient(t.hostport, id)
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
func (t *UDPTransporter) Start() {}
func (t *UDPTransporter) Stop()  {}

func (t *UDPTransporter) PrintInfo(inID, outID int) {
	fmt.Println("\033[1;33mUsage:\033[0m")
	fmt.Println()
}
