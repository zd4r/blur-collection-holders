package collection

import "sync"

type Map struct {
	sync.RWMutex
	m map[string]string
}

func NewMap() *Map {
	return &Map{
		m: make(map[string]string),
	}
}

func (m *Map) Set(key, val string) {
	m.Lock()
	defer m.Unlock()

	m.m[key] = val
}

func (m *Map) Get(key string) string {
	m.RLock()
	defer m.RUnlock()

	return m.m[key]
}

func (m *Map) Exists(key string) bool {
	m.RLock()
	defer m.RUnlock()

	_, ok := m.m[key]
	return ok
}
