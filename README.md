# Melrōse - programming of music melodies

[![Build Status](https://travis-ci.org/emicklei/melrose.png)](https://travis-ci.org/emicklei/melrose)
[![Go Report Card](https://goreportcard.com/badge/github.com/emicklei/melrose)](https://goreportcard.com/report/github.com/emicklei/melrose)
[![GoDoc](https://godoc.org/github.com/emicklei/melrose?status.svg)](https://pkg.go.dev/github.com/emicklei/melrose?tab=doc)


## Usage

This software can be used both as a library (Go package) or as part of the command line tool called `melrose` for live music coding.

### Use the tool

See [documentation](https://emicklei.github.io/melrose/) how to install and use `melrōse` for live performance.

### Use the Go Package

    package main

    import (
      "log"
      "time"

      m "github.com/emicklei/melrose"
      "github.com/emicklei/melrose/midi"
      "github.com/emicklei/melrose/op"
    )

    func main() {
      audio, err := midi.Open()
      if err != nil {
        log.Fatal(err)
      }
      defer audio.Close()

      f1 := m.MustParseSequence("C D E C")
      f2 := m.MustParseSequence("E F 2G")
      f3 := m.MustParseSequence("8G 8A 8G 8F E C")
      f4 := m.MustParseSequence("2C 2G3 2C 1=")
      r8 := op.Repeat{Target: []m.Sequenceable{m.Note("=")}, Times: m.On(8)}

      v1 := op.Join{Target: []m.Sequenceable{f1, f1, f2, f2, f3, f3, f4}}
      v2 := op.Join{Target: []m.Sequenceable{r8, v1}}
      v3 := op.Join{Target: []m.Sequenceable{r8, v2}}

      audio.Play(op.Join{Target: []m.Sequenceable{v1, v2, v3}}, 200, time.Now())
    }


Software is licensed under [Apache 2.0 license](LICENSE).
(c) 2014-2020 http://ernestmicklei.com 