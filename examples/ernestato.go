package main

import (
	"time"

	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

var Audio *audio.Device

func main() {
	Audio = new(audio.Device)
	Audio.Open()
	Audio.LoadSounds()
	defer Audio.Close()

	/*
		f2, _ := m.ParseSequence(`
			F#2 C#3 F#3 A3 C# F#`)
		Audio.Play(f2.Repeated(4), 800*time.Millisecond)

		f3, _ := m.ParseSequence(`
			F#2 (A3 C# F#)
			C#3 (A3 C# F#)
			F#3 (A3 C# F#)
			C#3 (A3 C# F#)`)
		Audio.Play(f3.Repeated(4), 800*time.Millisecond)

		time.Sleep(1 * time.Second)

		f4, _ := m.ParseSequence(`
			2F#2 (A3 C# F#) (A3 C# F#)
			2C#3 (A3 C# F#) (A3 C# F#)
			2F#3 (A3 C# F#) (A3 C# F#)
			2C#3 (A3 C# F#) (A3 C# F#)`)
		Audio.Play(f4.Repeated(4), 800*time.Millisecond)

		time.Sleep(1 * time.Second)
	*/

	/*
		f5, _ := ParseSequence(`
			G2 (B3 D G)
			D3 (B3 D G)
			G3 (B3 D G)
			D3 (B3 D G)`)
		Audio.Play(f5.Repeated(2), 800*time.Millisecond)
	*/

	b3dg := S("(B3 D G)")
	g2 := N("G2").S()
	d3 := N("D3").S()
	g3 := N("G3").S()
	Audio.Play(g2.Join(b3dg).Join(d3).Join(b3dg).Join(g3).Join(b3dg), 800*time.Millisecond)
}
