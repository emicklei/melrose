package main

import (
	"time"

	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audiolib"
)

func main() {
	Audio, _ := audiolib.Open()
	Audio.BeatsPerMinute(160)
	Audio.LoadSounds()
	defer Audio.Close()

	f2, _ := ParseSequence(`
			F#2 C#3 F#3 A3 C# F#`)
	Audio.Play(f2.Repeated(2))

	time.Sleep(1 * time.Second)

	f3, _ := ParseSequence(`
			F#2 (A3 C# F#)
			C#3 (A3 C# F#)
			F#3 (A3 C# F#)
			C#3 (A3 C# F#)`)
	Audio.Play(f3.Repeated(2))

	time.Sleep(1 * time.Second)

	f4, _ := ParseSequence(`
			2F#2 (A3 C# F#) (A3 C# F#)
			2C#3 (A3 C# F#) (A3 C# F#)
			2F#3 (A3 C# F#) (A3 C# F#)
			2C#3 (A3 C# F#) (A3 C# F#)`)
	Audio.Play(f4.Repeated(2))

	time.Sleep(1 * time.Second)

	f5, _ := ParseSequence(`
			G2 (B3 D G)
			D3 (B3 D G)
			G3 (B3 D G)
			D3 (B3 D G)`)
	Audio.Play(f5.Repeated(2))

	time.Sleep(1 * time.Second)

	b3dg := S("(B3 D G)")
	g2 := N("G2").S()
	d3 := N("D3").S()
	g3 := N("G3").S()
	Audio.Play(g2.Join(b3dg, d3, b3dg, g3, b3dg).S())
}
