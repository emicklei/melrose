package dsl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

// Snapshot is the object stored as JSON in a Save/Load operation.
type Snapshot struct {
	Author        string            `json:"author"`
	LastModified  time.Time         `json:"lastModified"`
	Syntax        string            `json:"syntax"`
	Variables     map[string]string `json:"variables"`
	Configuration map[string]string `json:"configuration"`
}

const defaultStorageFilename = "melrose.json"

// LoadMemoryFromDisk loads all variables by decoding JSON from a filename.
func LoadMemoryFromDisk(storage VariableStorage, args []string) notify.Message {
	filename := defaultStorageFilename
	if len(args) > 0 {
		filename = makeJSONFilename(args[0])
	}
	f, err := os.Open(filename)
	if err != nil {
		return notify.Errorf("unable to load:%v", err)
	}
	defer f.Close()

	snap := Snapshot{}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&snap); err != nil {
		return notify.Error(err)
	}

	if !IsCompatibleSyntax(snap.Syntax) {
		return notify.Errorf("syntax incompatible source detected, got %q want %q", snap.Syntax, Syntax)
	}

	eval := NewEvaluator(storage)
	toProcess := snap.Variables
	pass := 0
	for {
		if len(toProcess) == 0 {
			break
		}
		pass++
		if pass == 10 {
			break
		}
		toProcessNext := map[string]string{}
		for k, storex := range toProcess {
			v, err := eval.EvaluateExpression(storex)
			if err != nil {
				toProcessNext[k] = storex
				continue
			}
			storage.Put(k, v)
		}
		toProcess = toProcessNext
	}
	// check unresolveables
	if len(toProcess) > 0 {
		keys := []string{}
		for k := range toProcess {
			keys = append(keys, k)
		}
		return notify.Errorf("unable to evaluate variable(s) %v", keys)
	}

	// handle configuration
	if v, ok := snap.Configuration["bpm"]; ok {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return notify.Errorf("invalid value for bpm (beat-per-minute):%v", v)
		} else {
			// TODO do not like this dependency on a device
			if melrose.CurrentDevice() != nil {
				melrose.CurrentDevice().SetBeatsPerMinute(f)
			}
		}
	}
	return notify.Infof("loaded %d variables. use \":v\" to list them", len(snap.Variables))
}

// SaveMemoryToDisk saves all known variables in JSON to a filename.
func SaveMemoryToDisk(storage VariableStorage, args []string) notify.Message {
	filename := defaultStorageFilename
	if len(args) == 1 {
		filename = makeJSONFilename(args[0])
	}
	// make backup if exist
	if _, err := os.Stat(filename); err == nil {
		copyFile(filename, "backup_"+filename)
	}

	f, err := os.Create(filename)
	if err != nil {
		return notify.Errorf("unable to save:%v", err)
	}
	defer f.Close()

	storeMap := map[string]string{}
	for k, v := range storage.Variables() {
		if s, ok := v.(melrose.Storable); ok {
			storeMap[k] = s.Storex()
		} else {
			storeMap[k] = fmt.Sprintf("%v", v)
		}
	}

	snap := Snapshot{
		Author:        os.Getenv("USER"),
		LastModified:  time.Now(),
		Syntax:        Syntax,
		Variables:     storeMap,
		Configuration: map[string]string{},
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	if err := enc.Encode(snap); err != nil {
		return notify.Errorf("%v", err)
	}
	return notify.Infof("saved %d variables. use \":v\" to list them", len(storeMap))
}

func makeJSONFilename(entry string) string {
	if strings.HasSuffix(entry, ".json") {
		return entry
	}
	return entry + ".json"
}

func copyFile(from, to string) {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(to, input, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
