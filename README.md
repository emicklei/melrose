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

      "github.com/emicklei/melrose/dsl"
      "github.com/emicklei/melrose/midi"
    )

    func main() {
      audio, err := midi.Open()
      if err != nil {
        log.Fatal(err)
      }
      defer audio.Close()

      _, err = dsl.Run(audio, `
    bpm(120)
    play(sequence('C C# D D# E F G 2A- 2A#-- 2C5---'))
    `)
      if err != nil {
        log.Fatal(err)
      }
    }

Software is licensed under [Apache 2.0 license](LICENSE).
(c) 2014-2020 [ernestmicklei.com](http://ernestmicklei.com)
