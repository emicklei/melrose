package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/emicklei/melrose"
)

var memory = map[string]interface{}{}

type Var struct {
	Name string
}

func loadMemoryFromDisk() {
	f, err := os.Open(".melrose.image")
	if err != nil {
		printError(fmt.Sprintf("unable to load:%v", err))
		return
	}
	defer f.Close()

	storeMap := map[string]string{}
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&storeMap); err != nil {
		printError(err)
		return
	}

	// load into existing
	for k, s := range storeMap {
		v, err := eval(s)
		if err != nil {
			printError(fmt.Sprintf("unable to eval:%s = %s", k, s))
		} else {
			memory[k] = v
		}
	}

	printInfo(fmt.Sprintf("loaded %d variables. use \":v\" to list them", len(storeMap)))
}

func saveMemoryToDisk() {
	f, err := os.Create(".melrose.image")
	if err != nil {
		printError(fmt.Sprintf("unable to save:%v", err))
		return
	}
	defer f.Close()

	storeMap := map[string]string{}
	for k, v := range memory {
		if s, ok := v.(melrose.Storable); ok {
			storeMap[k] = s.Storex()
		} else {
			printError(fmt.Sprintf("cannot store %q:%T\n", k, v))
		}
	}

	enc := gob.NewEncoder(f)
	if err := enc.Encode(storeMap); err != nil {
		printError(err)
		return
	}
	printInfo(fmt.Sprintf("saved %d variables. use \":v\" to list them", len(storeMap)))
}
