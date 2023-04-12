package game

import (
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/ui/life"
)

type PlayerState struct {
	Program *tea.Program
	Id      int
	PosX    int
	PosY    int
	VelX    int
	VelY    int
	Paused  bool
	Color   int
	Placed  int
	Cells   int
}

type GameState int

type Lobby struct {
	players      map[int]*PlayerState
	playerColors [11]bool
	playersMutex sync.RWMutex
	playerCount  atomic.Int32
	incrementId  atomic.Int32
	board        [][]life.Cell
	boardMutex   sync.RWMutex
	ticker       *time.Ticker
	// state        GameState
}

const (
	PAUSED GameState = iota
	PLAYING
)

const MaxPlayers = 2
const MaxPlacedCells = 20
const drawRate = 20
const generationRate = 5
const drawsPerGeneration = drawRate / generationRate

const size = 100

func (l *Lobby) PlayerCount() int {
	return int(l.playerCount.Load())
}

func NewLobby() *Lobby {
	w, h := size, size

	return &Lobby{
		players:      make(map[int]*PlayerState),
		playerColors: [11]bool{true, false, false, false, false, false, false, false, false, false, false},
		board:        life.NewBoard(w, h),
		ticker:       time.NewTicker(time.Second / drawRate),
		// state:        PAUSED,
	}
}

func (l *Lobby) Run() {
	go func() {

		var prevUpdate time.Time
		iteration := 0

		for now := range l.ticker.C {
			iteration++

			if iteration == drawsPerGeneration {
				iteration = 0
				l.UpdateBoard()
			}

			l.Update(now.Sub(prevUpdate))

			prevUpdate = now
		}
	}()
}

type JoinLobbyMsg struct {
	Lobby       *Lobby
	PlayerState *PlayerState
	Id          int
	BoardWidth  int
	BoardHeight int
}

func (l *Lobby) Join(p *tea.Program) (int, bool) {
	l.playersMutex.Lock()
	defer l.playersMutex.Unlock()

	if l.playerCount.Load() == MaxPlayers {
		return 0, false
	}

	l.playerCount.Add(1)

	posX := rand.Intn(len(l.board))
	posY := rand.Intn(len(l.board))

	id := int(l.incrementId.Add(1))

	var color int
	for i := 0; i < 10; i++ {
		// id starts at 1, so in range [1, 10]
		// this cycles through all colors, even when players leave
		color = (i + id) % 11
		if !l.playerColors[color] {
			l.playerColors[color] = true
			break
		}
	}

	ps := PlayerState{
		Id:      id,
		Program: p,
		PosX:    posX,
		PosY:    posY,
		Paused:  true,
		Color:   color,
	}

	l.players[id] = &ps

	return id, true
}

func (l *Lobby) Leave(id int) {

	l.playersMutex.Lock()
	defer l.playersMutex.Unlock()
	l.boardMutex.Lock()
	defer l.boardMutex.Unlock()

	l.playerCount.Add(-1)

	l.playerColors[l.players[id].Color] = false
	delete(l.players, id)

	for y, row := range l.board {
		for x, cell := range row {
			if cell.Player == id {
				l.board[y][x].Player = life.DeadPlayer
				l.board[y][x].PausedPlayer = life.DeadPlayer
			}
		}
	}
}

func (l *Lobby) Unpause() {
	// update
}

func (l *Lobby) BoardSize() (int, int) {
	return len(l.board), len(l.board[0])
}

type UpdateBoardMsg struct{}

// type byPos []*ClientState

// func (s byPos) Len() int {
// 	return len(s)
// }
// func (s byPos) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }
// func (s byPos) Less(i, j int) bool {
// 	if s[i].PosX < s[j].PosX {
// 		return true
// 	}
// 	return s[i].PosY < s[j].PosY
// }

var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[0].cell))

func (l *Lobby) UpdateBoard() {

	// if l.state != PLAYING {
	// 	return
	// }

	l.boardMutex.Lock()
	defer l.boardMutex.Unlock()
	l.board = life.NextBoard(l.board)
}

func (l *Lobby) Update(delta time.Duration) {
	l.playersMutex.RLock()
	for _, player := range l.players {

		player.Program.Send(UpdateBoardMsg{})
	}
	l.playersMutex.RUnlock()

}

func (l *Lobby) ViewBoard(top, left, width, height int) string {

	// Arbitrary limits to avoid unreasonable terminal sizes
	// This already shows the board 4 times
	if width > size*2 {
		width = size * 2
	}
	if height > size*2 {
		height = size * 2
	}

	sb := strings.Builder{}

	boardWidth, boardHeight := l.BoardSize()

	l.playersMutex.RLock()
	defer l.playersMutex.RUnlock()
	l.boardMutex.RLock()
	defer l.boardMutex.RUnlock()

	for y := top; y < top+height; y++ {
		boundY := (y + boardHeight) % boardHeight

		deadCount := 0
		for x := left; x < left+width; x++ {
			boundX := (x + boardWidth) % boardWidth
			style := lipgloss.NewStyle()
			pixel := "  "

			cursor := false

			for _, player := range l.players {
				if boundY == player.PosY && boundX == player.PosX {
					cursor = true
					pixel = "[]"
					style = style.Foreground(lipgloss.Color(ColorTable[player.Color].cursor))

					break
				}
			}

			if l.board[boundY][boundX].Player == life.DeadPlayer && l.board[boundY][boundX].PausedPlayer == life.DeadPlayer && !cursor {
				deadCount++
				continue
			}
			sb.WriteString(deadStyle.Render(strings.Repeat("  ", deadCount)))
			deadCount = 0

			if !cursor {
				if l.board[boundY][boundX].Player != life.DeadPlayer {
					bc := l.players[l.board[boundY][boundX].Player].Color
					style = style.Background(lipgloss.Color(ColorTable[bc].cell))
				}
				if l.board[boundY][boundX].PausedPlayer != life.DeadPlayer {
					fc := l.players[l.board[boundY][boundX].PausedPlayer].Color
					style = style.Foreground(lipgloss.Color(ColorTable[fc].cell))
					pixel = "::"
				}
				// pixel = fmt.Sprint(l.board[boundY][boundX].Age)
			}
			sb.WriteString(style.Render(pixel))
		}

		if deadCount > 0 {
			sb.WriteString(deadStyle.Render(strings.Repeat("  ", deadCount)))
		}

		sb.WriteString("\n")
	}

	return sb.String()[:sb.Len()-1]
}

func (l *Lobby) Place(id int) {

	l.playersMutex.RLock()
	p, ok := l.players[id]
	l.playersMutex.RUnlock()

	// Maybe if player leaves but place() hasn't run yet?
	if !ok {
		return
	}
	// Can only place in pause mode
	if !p.Paused {
		return
	}

	l.boardMutex.Lock()
	defer l.boardMutex.Unlock()

	if l.board[p.PosY][p.PosX].PausedPlayer == life.DeadPlayer {
		if p.Placed >= MaxPlacedCells {
			return
		}
		l.board[p.PosY][p.PosX].PausedPlayer = p.Id
		p.Placed++

	} else if l.board[p.PosY][p.PosX].PausedPlayer == p.Color {
		l.board[p.PosY][p.PosX].PausedPlayer = 0
		p.Placed--
	}
}

func (l *Lobby) TogglePause(id int) {
	l.playersMutex.RLock()
	p, ok := l.players[id]
	l.playersMutex.RUnlock()

	// Maybe if player leaves unexpectedly
	if !ok {
		return
	}

	l.boardMutex.Lock()
	defer l.boardMutex.Unlock()

	p.Paused = !p.Paused

	if !p.Paused {
		// valid := true

		for _, row := range l.board {
			for _, cell := range row {
				if cell.PausedPlayer == id && cell.Player != life.DeadPlayer {
					// valid = false
					return
				}
			}
		}
		for y, row := range l.board {
			for x, cell := range row {
				if cell.PausedPlayer == id {
					l.board[y][x].Player = cell.PausedPlayer
				}
			}
		}

	} else {
		for y, row := range l.board {
			for x, cell := range row {
				if cell.Player == id {
					l.board[y][x].Player = life.DeadPlayer
				}
			}
		}
		// Don't need a lock, b/c
		p.Placed = 0
	}

}

func (l *Lobby) GetPlayer(id int) *PlayerState {
	// TODO i can't tell if this mutex is rlock is necessary
	l.playersMutex.RLock()
	defer l.playersMutex.RUnlock()
	return l.players[id]
}
