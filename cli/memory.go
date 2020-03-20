package main

import (
	"encoding/gob"
	"log"
	"os"
)

var memory = map[string]interface{}{}

type Var struct {
	Name string
}

func loadMemoryFromDisk() {

}

func saveMemoryToDisk() {
	f, _ := os.Create(".melrose.image")
	defer f.Close()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(memory); err != nil {
		log.Fatal(err)
	}
}
