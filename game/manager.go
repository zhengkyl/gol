package game

import (
	"fmt"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/zhengkyl/gol/game/life"
)

type programState struct {
	program *tea.Program
	lobbyId int
}

const (
	lobbyIdMenu = -2
	lobbyIdSolo = -1
)

type Manager struct {
	lobbies      map[int]*Lobby
	lobbiesMutex sync.RWMutex
	lobbyId      int
	players      map[int]programState
	playersMutex sync.RWMutex
	playerId     int
}

func NewManager() *Manager {
	return &Manager{
		lobbies: make(map[int]*Lobby),
		players: make(map[int]programState),
		// lobbyId: ,
	}
}

func (gm *Manager) NewLobby() *Lobby {
	w, h := defaultWidth, defaultHeight

	l := &Lobby{
		players:      make(map[int]*PlayerState),
		playerColors: [11]bool{true, false, false, false, false, false, false, false, false, false, false},
		board:        life.NewBoard(w, h),
		ticker:       time.NewTicker(time.Second / drawRate),
		name:         petname.Generate(2, " "),
	}

	gm.lobbiesMutex.Lock()
	gm.lobbyId++
	l.id = gm.lobbyId
	gm.lobbies[l.id] = l
	gm.lobbiesMutex.Unlock()

	gm.broadcastLobbyInfos()

	return l
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

type LobbyInfo struct {
	PlayerCount int
	MaxPlayers  int
	Name        string
	Id          int
}

func (gm *Manager) broadcastLobbyInfos() {
	infos := gm.LobbyInfos()

	gm.playersMutex.RLock()
	for _, ps := range gm.players {
		if ps.lobbyId == lobbyIdMenu {
			ps.program.Send(infos)
		}
	}
	gm.playersMutex.RUnlock()
}

func (gm *Manager) LobbyInfos() []LobbyInfo {
	infos := make([]LobbyInfo, 0)
	gm.lobbiesMutex.RLock()
	for _, l := range gm.lobbies {
		infos = append(infos, LobbyInfo{
			PlayerCount: l.playerCount,
			MaxPlayers:  MaxPlayers,
			Name:        l.name,
			Id:          l.id,
		})
	}
	gm.lobbiesMutex.RUnlock()
	return infos
}

func (gm *Manager) Connect(p *tea.Program) int {
	gm.playersMutex.Lock()
	defer gm.playersMutex.Unlock()

	gm.playerId++

	gm.players[gm.playerId] = programState{
		program: p,
		lobbyId: lobbyIdMenu,
	}

	return gm.playerId
}

func (gm *Manager) Disconnect(playerId int) {
	gm.playersMutex.Lock()
	defer gm.playersMutex.Unlock()

	state, ok := gm.players[playerId]
	if !ok {
		return // maybe disconnect before connect? idk if possible
	}

	if state.lobbyId >= 0 {
		gm.removeFromLobby(state.lobbyId, playerId)
	}

	delete(gm.players, playerId)
}

type JoinSuccessMsg struct {
	Lobby       *Lobby
	PlayerState *PlayerState
	Id          int
	BoardWidth  int
	BoardHeight int
}
type JoinFailMsg struct {
	err string
}

func (gm *Manager) JoinLobby(lobbyId int, playerId int) tea.Msg {
	gm.lobbiesMutex.RLock()
	lobby, ok := gm.lobbies[lobbyId]
	gm.lobbiesMutex.RUnlock()

	if !ok {
		return JoinFailMsg{fmt.Sprintf("Lobby with id=%v does not exist", lobbyId)}
	}

	gm.playersMutex.Lock()
	defer gm.playersMutex.Unlock()

	program := gm.players[playerId].program
	err := lobby.Join(playerId, program)
	if err != nil {
		return JoinFailMsg{err.Error()}
	}

	gm.players[playerId] = programState{program: program, lobbyId: lobbyId}

	return JoinSuccessMsg{}
}

func (gm *Manager) removeFromLobby(lobbyId, playerId int) {
	gm.lobbiesMutex.RLock()
	lobby, ok := gm.lobbies[lobbyId]

	// Lobby is gone, so we good
	if !ok {
		gm.lobbiesMutex.RUnlock()
		return
	}

	lobby.Leave(playerId)
	count := lobby.playerCount
	gm.lobbiesMutex.RUnlock()

	if count == 0 {
		gm.lobbiesMutex.Lock()
		delete(gm.lobbies, lobbyId)
		gm.lobbiesMutex.Unlock()

		gm.broadcastLobbyInfos()
	}

}

func (gm *Manager) LeaveLobby(playerId int) {
	gm.playersMutex.Lock()
	state := gm.players[playerId]
	gm.players[playerId] = programState{lobbyId: lobbyIdMenu}
	gm.playersMutex.Unlock()

	gm.removeFromLobby(state.lobbyId, playerId)
}
