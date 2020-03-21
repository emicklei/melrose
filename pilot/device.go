package pilot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	m "github.com/emicklei/melrose"
)

// https://github.com/hundredrabbits/Pilot/blob/master/desktop/sources/scripts/mixer.js
// p.Send("1OSCsisq")
// p.Send("reset")
// p.Send("rosc")

type Pilot struct {
	conn    net.Conn
	channel int
}

func Open() (*Pilot, error) {
	conn, err := net.Dial("udp4", ":49161")
	if err != nil {
		return nil, err
	}
	return &Pilot{conn: conn, channel: 0}, nil
}
func (p *Pilot) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Pilot) Play(seq m.Sequence) {
	seq.NotesDo(func(n m.Note) {
		cmd := noteToCmd(p.channel, n)
		p.Send(cmd)
	})
}

func (p *Pilot) Send(cmd string) {
	log.Println(cmd)
	_, err := fmt.Fprintf(p.conn, cmd)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(200 * time.Millisecond)
}

func noteToCmd(ch int, n m.Note) string {
	var b bytes.Buffer
	io.WriteString(&b, strconv.Itoa(ch)) // channel
	io.WriteString(&b, strconv.Itoa(n.Octave))
	io.WriteString(&b, n.Name)
	io.WriteString(&b, "f")
	return b.String()
}
