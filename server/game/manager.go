package game

import (
	"sync"
)

type Manager struct {
	lobbies      map[*Lobby]struct{}
	lobbiesMutex sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		lobbies: make(map[*Lobby]struct{}),
	}
}

func (gm *Manager) FindLobby() *Lobby {
	gm.lobbiesMutex.RLock()
	for g := range gm.lobbies {
		if g.PlayerCount() >= MaxPlayers {
			continue
		}
		gm.lobbiesMutex.RUnlock()
		return g
	}
	gm.lobbiesMutex.RUnlock()

	g := NewLobby()
	g.Run()

	gm.lobbiesMutex.Lock()
	gm.lobbies[g] = struct{}{}
	gm.lobbiesMutex.Unlock()

	return g
}

func (gm *Manager) Games(g *Lobby) {

}

func (gm *Manager) EndGame(g *Lobby) {
	delete(gm.lobbies, g)
}
