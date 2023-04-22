package game

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/zhengkyl/gol/game/life"
)

type programState struct {
	lobbyId  int
	playerId int
}

const (
	lobbyIdMenu = -2
	lobbyIdSolo = -1
)

type Manager struct {
	lobbies       map[int]*Lobby
	lobbiesMutex  sync.RWMutex
	lobbyId       atomic.Int32
	programs      map[*tea.Program]programState
	programsMutex sync.RWMutex
	programCount  int
}

func NewManager() *Manager {
	return &Manager{
		lobbies:  make(map[int]*Lobby),
		programs: make(map[*tea.Program]programState),
		// lobbyId: ,
	}
}

func (gm *Manager) NewLobby() *Lobby {
	w, h := defaultWidth, defaultHeight

	gm.lobbiesMutex.Lock()

	return &Lobby{
		players:      make(map[int]*PlayerState),
		playerColors: [11]bool{true, false, false, false, false, false, false, false, false, false, false},
		board:        life.NewBoard(w, h),
		ticker:       time.NewTicker(time.Second / drawRate),
		name:         petname.Generate(2, " "),
		id:           int(gm.lobbyId.Add(1)),
	}
}

// TODO maybe auto find lobby button?
// func (gm *Manager) FindLobby() *Lobby {
// 	gm.lobbiesMutex.RLock()
// 	for _, g := range gm.lobbies {
// 		if g.PlayerCount() >= MaxPlayers {
// 			continue
// 		}
// 		gm.lobbiesMutex.RUnlock()
// 		return g
// 	}
// 	gm.lobbiesMutex.RUnlock()

// 	g := gm.NewLobby()
// 	g.Run()

// 	return g
// }

type LobbyStatus struct {
	PlayerCount int
	MaxPlayers  int
	Name        string
	Id          int
}

// TODO send msg to all programs on menu when jion/leave
// if lobby is destroyed, pointer is invalid
func (gm *Manager) LobbyStatuses() []LobbyStatus {
	statuses := make([]LobbyStatus, 0)
	gm.lobbiesMutex.RLock()
	for _, g := range gm.lobbies {
		statuses = append(statuses, LobbyStatus{
			PlayerCount: int(g.playerCount.Load()),
			MaxPlayers:  MaxPlayers,
			Name:        g.name,
			Id:          g.id,
		})
	}
	gm.lobbiesMutex.RUnlock()
	return statuses
}

func (gm *Manager) Connect(p *tea.Program) {
	gm.programsMutex.Lock()
	defer gm.programsMutex.Unlock()

	gm.programs[p] = programState{
		lobbyId: lobbyIdMenu,
	}
	gm.programCount++

}

func (gm *Manager) Disconnect(p *tea.Program) {
	gm.programsMutex.Lock()
	defer gm.programsMutex.Unlock()

	state, ok := gm.programs[p]
	if !ok {
		return // maybe disconnect before connect? idk if possible
	}

	if state.lobbyId >= 0 {
		gm.LeaveLobby(state.lobbyId, state.playerId)
	}

	delete(gm.programs, p)
	gm.programCount--
}

func (gm *Manager) JoinLobby(lobbyId int, p *tea.Program) (int, error) {
	gm.lobbiesMutex.RLock()

	lobby, ok := gm.lobbies[lobbyId]

	if !ok {
		return 0, fmt.Errorf("Lobby with id=%v does not exist", lobbyId)
	}

	playerId, err := lobby.Join(p)

	gm.lobbiesMutex.RUnlock()

	if err != nil {
		return 0, err
	}

	gm.programsMutex.Lock()
	gm.programs[p] = programState{lobbyId: lobbyId, playerId: playerId}
	gm.programsMutex.Unlock()

	return playerId, nil
}

func (gm *Manager) LeaveLobby(lobbyId, playerId int) {
	gm.lobbiesMutex.RLock()
	lobby, ok := gm.lobbies[lobbyId]

	// Lobby is gone, so we good
	if !ok {
		gm.lobbiesMutex.RUnlock()
		return
	}

	lobby.Leave(playerId)
	count := lobby.playerCount.Load()
	gm.lobbiesMutex.RUnlock()

	if count == 0 {
		gm.lobbiesMutex.Lock()
		delete(gm.lobbies, lobbyId)
		gm.lobbiesMutex.Unlock()
	}
}

func (gm *Manager) EscToMenu(p *tea.Program) {
	gm.programsMutex.Lock()
	state := gm.programs[p]
	gm.programs[p] = programState{lobbyId: lobbyIdMenu}
	gm.programsMutex.Unlock()

	gm.LeaveLobby(state.lobbyId, state.playerId)
}
