package game

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

type PlayerMap struct {
	mu sync.Mutex
	v  map[*tea.Program]*ClientState
}

func (m *PlayerMap) Set(p *tea.Program, cs *ClientState) {
	m.mu.Lock()
	m.v[p] = cs
	m.mu.Unlock()
}

func (m *PlayerMap) Delete(p *tea.Program) {
	m.mu.Lock()
	delete(m.v, p)
	m.mu.Unlock()
}

func (m *PlayerMap) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.v)
}

func (m *PlayerMap) Entries() ([]*tea.Program, []*ClientState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := make([]*tea.Program, 0, len(m.v))
	vals := make([]*ClientState, 0, len(m.v))
	for p, cs := range m.v {
		keys = append(keys, p)
		vals = append(vals, cs)
	}

	return keys, vals
}
