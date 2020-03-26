package pilot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/emicklei/melrose"
)

// https://github.com/hundredrabbits/Pilot/blob/master/desktop/sources/scripts/mixer.js
// p.Send("1OSCsisq")
// p.Send("reset")
// p.Send("rosc")

type Pilot struct {
	enabled bool
	conn    net.Conn
	channel int
	bpm     float64
}

func Open() (*Pilot, error) {
	conn, err := net.Dial("udp4", ":49161")
	if err != nil {
		fmt.Println("pilot connect error", err)
		return nil, err
	}
	return &Pilot{conn: conn, channel: 0, enabled: true, bpm: 0}, nil
}
func (p *Pilot) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Pilot) Play(seq melrose.Sequence, echo bool) {
	if !p.enabled {
		return
	}
	wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/p.bpm))) * time.Millisecond
	for _, eachGroup := range seq.Notes {
		if len(eachGroup) == 1 {
			print(eachGroup[0])
			p.PlayNote(eachGroup[0], wholeNoteDuration)
			continue
		}
		// grouped
		print(eachGroup)
		multiCmd := ""
		for i, note := range eachGroup {
			if i > 0 {
				multiCmd += ";"
			}
			cmd := noteToCmd(p.channel, note)
			multiCmd += cmd

		}
		p.send(multiCmd)
		actualDuration := time.Duration(float32(wholeNoteDuration) * eachGroup[0].DurationFactor()) // TODO first note
		time.Sleep(actualDuration)
	}
	fmt.Println()
}

func (p *Pilot) PlayNote(note melrose.Note, wholeNoteDuration time.Duration) {
	actualDuration := time.Duration(float32(wholeNoteDuration) * note.DurationFactor())
	if note.IsRest() {
		time.Sleep(actualDuration)
		return
	}
	cmd := noteToCmd(p.channel, note)
	// note on
	p.send(cmd)
	time.Sleep(actualDuration)
	// no note off
}

func (p *Pilot) send(cmd string) {
	//log.Println("pilot debug:", cmd)
	_, err := fmt.Fprintf(p.conn, cmd)
	if err != nil {
		p.enabled = false
		log.Println("pilot send error", err)
	}
}

func noteToCmd(ch int, n melrose.Note) string {
	var b bytes.Buffer
	io.WriteString(&b, strconv.Itoa(ch)) // channel
	io.WriteString(&b, strconv.Itoa(n.Octave))
	io.WriteString(&b, n.Name)
	io.WriteString(&b, "f")
	return b.String()
}

// BeatsPerMinute (BPM) ; beats each the length of a quarter note per minute.
func (p *Pilot) SetBeatsPerMinute(bpm float64) {
	if bpm <= 0 {
		return
	}
	p.send(fmt.Sprintf("bpm%d", int(bpm)))
	p.bpm = bpm
}

func (p *Pilot) BeatsPerMinute() float64 {
	return p.bpm
}

// 34 = purple
func print(arg interface{}) {
	fmt.Printf("\033[2;34m" + fmt.Sprintf("%v ", arg) + "\033[0m")
}
