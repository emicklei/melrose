// +build udp

package transport

import "flag"

var transportUDP = flag.Int("udp", 0, "if set to port > 0 then use UDP transport")

func init() { UseUDPTransport(9000) }
