package address

import "sync"

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
