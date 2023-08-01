package game

import (
	"fmt"
	"sort"
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
	}
}

func (gm *Manager) CreateLobby() int {
	w, h := defaultWidth, defaultHeight

	l := &Lobby{
		players:      make(map[int]*PlayerState),
		playerColors: [11]bool{true, false, false, false, false, false, false, false, false, false, false},
		board:        life.NewBoard(w, h),
		ticker:       time.NewTicker(time.Second / drawRate),
		name:         petname.Generate(2, "-"),
	}

	gm.lobbiesMutex.Lock()
	gm.lobbyId++
	l.id = gm.lobbyId
	gm.lobbies[l.id] = l
	gm.lobbiesMutex.Unlock()

	gm.BroadcastLobbyInfos()

	l.Run()

	return l.id
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

type SoloGameMsg struct{}
type LobbyInfoList []LobbyInfo

type LobbyInfo struct {
	PlayerCount int
	MaxPlayers  int
	Name        string
	Id          int
}

func (gm *Manager) BroadcastLobbyInfos() {
	infos := gm.LobbyInfos()

	gm.playersMutex.RLock()
	for _, ps := range gm.players {
		if ps.lobbyId == lobbyIdMenu {
			go func(p *tea.Program) {
				p.Send(infos)
			}(ps.program)
		}
	}
	gm.playersMutex.RUnlock()
}

func (gm *Manager) Debug() string {

	gm.playersMutex.RLock()
	defer gm.playersMutex.RUnlock()
	return fmt.Sprint(gm.players)
}

type byId []LobbyInfo

func (s byId) Len() int {
	return len(s)
}
func (s byId) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byId) Less(i, j int) bool {
	return s[i].Id < s[j].Id
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

	sort.Sort(byId(infos))
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
	state, ok := gm.players[playerId]
	gm.playersMutex.Unlock()

	if !ok {
		return // maybe disconnect before connect? idk if possible
	}

	if state.lobbyId >= 0 {
		gm.removeFromLobby(state.lobbyId, playerId)
	}

	delete(gm.players, playerId)
}

func (gm *Manager) LeaveLobby(playerId int) {
	gm.playersMutex.Lock()
	state, ok := gm.players[playerId]
	if ok {
		gm.players[playerId] = programState{lobbyId: lobbyIdMenu, program: state.program}
	}
	gm.playersMutex.Unlock()

	if !ok {
		return
	}

	if state.lobbyId >= 0 {
		gm.removeFromLobby(state.lobbyId, playerId)
	}
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

	program := gm.players[playerId].program
	ps, err := lobby.Join(playerId, program)
	if err != nil {
		gm.playersMutex.Unlock()
		return JoinFailMsg{err.Error()}
	}

	gm.players[playerId] = programState{program: program, lobbyId: lobbyId}

	gm.playersMutex.Unlock()
	bw, bh := lobby.BoardSize()

	gm.BroadcastLobbyInfos()

	return JoinSuccessMsg{
		Lobby:       lobby,
		PlayerState: ps,
		Id:          lobbyId,
		BoardWidth:  bw,
		BoardHeight: bh,
	}
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
	}

	gm.BroadcastLobbyInfos()

}
