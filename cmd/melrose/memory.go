package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/emicklei/melrose"
)

var memory = map[string]interface{}{}

type Var struct {
	Name string
}

func (v Var) Value() interface{} {
	return memory[v.Name]
}

func (v Var) Storex() string {
	return fmt.Sprintf(`var(%s)`, v.Name)
}

func (v Var) String() string {
	return fmt.Sprintf(`%s:%T`, v.Name, v.Value())
}

func (v Var) S() melrose.Sequence {
	if s, ok := v.Value().(melrose.Sequenceable); ok {
		return s.S()
	}
	return melrose.Sequence{}
}

func variableNameFor(value interface{}) string {
	for k, v := range memory {
		if reflect.DeepEqual(value, v) {
			return k
		}
	}
	return "" // not found
}

func listVariables() {
	keys := []string{}
	width := 0
	for k, _ := range memory {
		if len(k) > width {
			width = len(k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := memory[k]
		if s, ok := v.(melrose.Storable); ok {
			fmt.Printf("%s = %s\n", strings.Repeat(" ", width-len(k))+k, s.Storex())
		} else {
			fmt.Printf("%s = (%T) %v\n", strings.Repeat(" ", width-len(k))+k, v, v)
		}
	}
}

const defaultStorageFilename = "melrose.json"

func loadMemoryFromDisk() {
	f, err := os.Open(defaultStorageFilename)
	if err != nil {
		printError(fmt.Sprintf("unable to load:%v", err))
		return
	}
	defer f.Close()

	storeMap := map[string]string{}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&storeMap); err != nil {
		printError(err)
		return
	}

	// if var is used and its value is not known we do a second pass. TODO workaround fix
	secondsPass := map[string]string{}
	// load into existing
	for k, s := range storeMap {
		v, err := eval(s)
		if err != nil {
			secondsPass[k] = s
		} else {
			memory[k] = v
		}
	}
	for k, s := range secondsPass {
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
	f, err := os.Create(defaultStorageFilename)
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

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	if err := enc.Encode(storeMap); err != nil {
		printError(err)
		return
	}
	printInfo(fmt.Sprintf("saved %d variables. use \":v\" to list them", len(storeMap)))
}
