package dsl

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

type variable struct {
	Name  string
	store *VariableStore
}

func (v variable) Storex() string {
	return v.Name // fmt.Sprintf(`var(%s)`, v.Name)
}

func (v variable) S() melrose.Sequence {
	m, ok := v.store.Get(v.Name)
	if !ok {
		return melrose.Sequence{}
	}
	if s, ok := m.(melrose.Sequenceable); ok {
		return s.S()
	}
	return melrose.Sequence{}
}

// VariableStore is an in-memory storage of values by name.
// Access to this store is go-routine safe.
type VariableStore struct {
	mutex     sync.RWMutex
	variables map[string]interface{}
}

// NewVariableStore returns a new
func NewVariableStore() *VariableStore {
	return &VariableStore{
		variables: map[string]interface{}{},
	}
}

// NameFor finds the entry for a value and returns its (first) associated name
func (v *VariableStore) NameFor(value interface{}) string {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	for k, v := range v.variables {
		if reflect.DeepEqual(value, v) {
			// the first found
			return k
		}
	}
	return "" // not found
}

// Get returns a value found by the key. Inspect the second return value of presence.
func (v *VariableStore) Get(key string) (interface{}, bool) {
	v.mutex.RLock()
	e, ok := v.variables[key]
	v.mutex.RUnlock()
	return e, ok
}

// Put stores a value by the key. Overwrites any existing value.
func (v *VariableStore) Put(key string, value interface{}) {
	v.mutex.Lock()
	v.variables[key] = value
	v.mutex.Unlock()
}

// Delete removes a stored value by the key. Ignores if the key is not found.
func (v *VariableStore) Delete(key string) {
	v.mutex.Lock()
	delete(v.variables, key)
	v.mutex.Unlock()
}

// Variables returns a copy of all stores variables.
func (v *VariableStore) Variables() map[string]interface{} {
	v.mutex.RLock()
	copy := map[string]interface{}{}
	for k, v := range v.variables {
		copy[k] = v
	}
	v.mutex.RUnlock()
	return copy
}

// ListVariables prints a list of sorted key=value pairs.
func (v *VariableStore) ListVariables(entry string) notify.Message {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	keys := []string{}
	width := 0
	for k, _ := range v.variables {
		if len(k) > width {
			width = len(k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := v.variables[k]
		if s, ok := v.(melrose.Storable); ok {
			fmt.Printf("%s = %s\n", strings.Repeat(" ", width-len(k))+k, s.Storex())
		} else {
			fmt.Printf("%s = (%T) %v\n", strings.Repeat(" ", width-len(k))+k, v, v)
		}
	}
	return nil
}

// Snapshot is the object stored as JSON in a Save/Load operation.
type Snapshot struct {
	Variables     map[string]string `json:"variables"`
	Configuration map[string]string `json:"configuration"`
}

const defaultStorageFilename = "melrose.json"

// LoadMemoryFromDisk loads all variables by decoding JSON from a filename.
func (s *VariableStore) LoadMemoryFromDisk(entry string) notify.Message {
	f, err := os.Open(defaultStorageFilename)
	if err != nil {
		return notify.Errorf("unable to load:%v", err)
	}
	defer f.Close()

	snap := Snapshot{}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&snap); err != nil {
		return notify.Error(err)
	}

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
			v, err := Evaluate(s, storex)
			if err != nil {
				toProcessNext[k] = storex
				continue
			}
			if r, ok := v.(FunctionResult); ok {
				s.Put(k, r.Result)
			}
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
			melrose.CurrentDevice().SetBeatsPerMinute(f)
		}
	}
	return notify.Infof("loaded %d variables. use \":v\" to list them", len(snap.Variables))
}

// SaveMemoryToDisk saves all known variables in JSON to a filename.
func (s *VariableStore) SaveMemoryToDisk(entry string) notify.Message {
	f, err := os.Create(defaultStorageFilename)
	if err != nil {
		return notify.Errorf("unable to save:%v", err)
	}
	defer f.Close()

	storeMap := map[string]string{}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for k, v := range s.variables {
		if s, ok := v.(melrose.Storable); ok {
			storeMap[k] = s.Storex()
		} else {
			return notify.Errorf("cannot store %q:%T\n", k, v)
		}
	}

	snap := Snapshot{
		Variables: storeMap,
		Configuration: map[string]string{
			"bpm": fmt.Sprintf("%v", melrose.CurrentDevice().BeatsPerMinute()),
		},
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	if err := enc.Encode(snap); err != nil {
		return notify.Errorf("%v", err)
	}
	return notify.Infof("saved %d variables. use \":v\" to list them", len(storeMap))
}
