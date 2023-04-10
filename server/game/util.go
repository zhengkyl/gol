package game

import (
	"sync"
)

type PlayerMap struct {
	mu        sync.Mutex
	playerMap map[int]*PlayerState
	counter   int
}

func (m *PlayerMap) Add(cs *PlayerState) int {
	m.mu.Lock()
	m.counter++
	m.playerMap[m.counter] = cs
	m.mu.Unlock()
	return m.counter
}

func (m *PlayerMap) Delete(id int) {
	m.mu.Lock()
	delete(m.playerMap, id)
	m.mu.Unlock()
}

func (m *PlayerMap) Len() int {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return len(m.playerMap)
}

func (m *PlayerMap) Get(id int) *PlayerState {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return m.playerMap[id]
}

func (m *PlayerMap) Entries() ([]int, []*PlayerState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := make([]int, 0, len(m.playerMap))
	vals := make([]*PlayerState, 0, len(m.playerMap))
	for id, ps := range m.playerMap {
		keys = append(keys, id)
		vals = append(vals, ps)
	}

	return keys, vals
}
