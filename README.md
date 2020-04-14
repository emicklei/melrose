# Melrose - programming of music melodies

[![Go Report Card](https://goreportcard.com/badge/github.com/emicklei/melrose)](https://goreportcard.com/report/github.com/emicklei/melrose)
[![GoDoc](https://godoc.org/github.com/emicklei/melrose?status.svg)](https://pkg.go.dev/github.com/emicklei/melrose?tab=doc)


##

## Usage

This software can be used both as a library (Go package) or as part of the command line tool called `melrose` for live music coding.

### Use the tool

See [documentation](https://emicklei.github.io/melrose/) how to install and use `melrose` for live performance.

### Use the Go Package

    package main

    import (
      "log"

      "github.com/emicklei/melrose/m"
      "github.com/emicklei/melrose/midi"
    )

    func main() {
      audio, err := midi.Open()
      if err != nil {
        log.Fatal(err)
      }
      audio.SetBeatsPerMinute(200)
      defer audio.Close()

      f1 := m.Sequence("C D E C")
      f2 := m.Sequence("E F ½G")
      f3 := m.Sequence("8G 8A 8G 8F E C")
      f4 := m.Sequence("2C 2G3 2C 1=")
      r8 := m.Repeat(8, m.Note("="))

      v1 := m.Join(f1, f1, f2, f2, f3, f3, f4)
      v2 := m.Join(r8, v1)
      v3 := m.Join(r8, v2)

      m.Go(audio, v1, v2, v3)
    }


Software is licensed under [Apache 2.0 license](LICENSE).
(c) 2014-2020 http://ernestmicklei.com 