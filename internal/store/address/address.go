package address

import (
	"sort"
	"sync"
)

type Map struct {
	sync.RWMutex
	m map[string]struct{}
}

func NewMap() *Map {
	return &Map{
		m: make(map[string]struct{}),
	}
}

func (m *Map) Set(address string) {
	m.Lock()
	defer m.Unlock()

	m.m[address] = struct{}{}
}

func (m *Map) Exists(address string) bool {
	m.RLock()
	defer m.RUnlock()

	_, ok := m.m[address]

	return ok
}

func (m *Map) Clear() {
	m.RLock()
	defer m.RUnlock()

	clear(m.m)
}

func (m *Map) GetAll() []string {
	m.RLock()
	defer m.RUnlock()

	keys := make([]string, 0, len(m.m))

	for key, _ := range m.m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
