package dsl

import (
	"fmt"
	"github.com/emicklei/melrose/core"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/emicklei/melrose/notify"
)

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

func (v *VariableStore) getVariable(name string) variable {
	return variable{Name: name, store: v}
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
func ListVariables(storage core.VariableStorage, args []string) notify.Message {
	keys := []string{}
	width := 0
	variables := storage.Variables()
	for k, _ := range variables {
		// if filtering is wanted
		if len(args) == 1 && !strings.HasPrefix(k, args[0]) {
			continue
		}
		if len(k) > width {
			width = len(k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := variables[k]
		if s, ok := v.(core.Storable); ok {
			fmt.Printf("%s = %s\n", strings.Repeat(" ", width-len(k))+k, s.Storex())
		} else {
			fmt.Printf("%s = (%T) %v\n", strings.Repeat(" ", width-len(k))+k, v, v)
		}
	}
	return nil
}
