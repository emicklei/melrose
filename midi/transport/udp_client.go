package transport

import (
	"fmt"
	"net"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi/io"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/tre"
)

type routerClient struct {
	connection net.Conn
}

func newRouterClient(port, id int) (MIDIOut, error) {
	con, err := net.Dial("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, tre.New(err, "net.Dial", "port", port)
	}
	return routerClient{connection: con}, nil
}

func (r routerClient) Abort() error { return nil }

func (r routerClient) Close() error {
	return r.connection.Close()
}

func (r routerClient) WriteShort(status int64, data1 int64, data2 int64) error {
	if core.IsDebug() {
		notify.Debugf("transport.RouterClient.WriteShort %d,%d,%d", status, data1, data2)
	}
	return io.WriteMessage(status, data1, data2, r.connection)
}
