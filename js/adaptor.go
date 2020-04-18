package js

import (
	"sync"

	"github.com/dop251/goja"
	"github.com/emicklei/melrose/dsl"
)

// storageAdaptor implements dsl.VariableStorage
type storageAdaptor struct {
	vm     *goja.Runtime
	mutex  *sync.RWMutex
	mirror dsl.VariableStorage
}

func newAdaptorOn(vm *goja.Runtime, store dsl.VariableStorage) *storageAdaptor {
	return &storageAdaptor{
		vm:     vm,
		mutex:  new(sync.RWMutex),
		mirror: store,
	}
}

func (s *storageAdaptor) NameFor(value interface{}) string {
	// mirror is go-safe
	return s.mirror.NameFor(value)
}

func (s *storageAdaptor) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v := s.vm.Get(key).Export()
	return v, v == nil
}

func (s *storageAdaptor) Put(key string, value interface{}) {
	s.mutex.Lock()
	s.vm.Set(key, value)
	s.mutex.Unlock()
	// mirror is go-safe
	s.mirror.Put(key, value)
}

func (s *storageAdaptor) Delete(key string) {
	s.mutex.Lock()
	s.vm.Set(key, nil)
	s.mutex.Unlock()
	// mirror is go-safe
	s.mirror.Delete(key)
}

func (s *storageAdaptor) Variables() map[string]interface{} {
	// mirror is go-safe
	return s.mirror.Variables()
}
