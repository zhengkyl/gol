package game

import (
	"fmt"
	"math/rand"
	"sort"
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
}

const MaxPlayers = 10
const MaxPlacedCells = 40
const drawRate = 20
const generationRate = 5
const drawsPerGeneration = drawRate / generationRate

const defaultWidth = 160
const defaultHeight = 90

func NewLobby() *Lobby {
	w, h := defaultWidth, defaultHeight

	return &Lobby{
		players:      make(map[int]*PlayerState),
		playerColors: [11]bool{true, false, false, false, false, false, false, false, false, false, false},
		board:        life.NewBoard(w, h),
		ticker:       time.NewTicker(time.Second / drawRate),
	}
}

func (l *Lobby) PlayerCount() int {
	return int(l.playerCount.Load())
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
	return len(l.board[0]), len(l.board)
}

type UpdateBoardMsg struct{}

type byCells []*PlayerState

func (s byCells) Len() int {
	return len(s)
}
func (s byCells) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCells) Less(i, j int) bool {
	if s[i].Cells < s[j].Cells {
		return true
	}
	if s[i].Cells > s[j].Cells {
		return false
	}
	return s[i].Color < s[j].Color
}

var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[0].Cell))

func (l *Lobby) UpdateBoard() {

	l.boardMutex.Lock()
	l.board = life.NextBoard(l.board)
	l.boardMutex.Unlock()

	l.playersMutex.Lock()
	for _, ps := range l.players {
		ps.Cells = 0
	}

	for _, row := range l.board {
		for _, cell := range row {
			if cell.Player != life.DeadPlayer {
				l.players[cell.Player].Cells++
			}
		}
	}
	l.playersMutex.Unlock()
}

func (l *Lobby) Update(delta time.Duration) {
	l.playersMutex.RLock()
	for _, player := range l.players {

		player.Program.Send(UpdateBoardMsg{})
	}
	l.playersMutex.RUnlock()

}

func (l *Lobby) Scoreboard() string {
	l.playersMutex.RLock()
	defer l.playersMutex.RUnlock()

	sb := strings.Builder{}
	var ps []*PlayerState

	for _, p := range l.players {
		ps = append(ps, p)
	}
	sort.Sort(sort.Reverse(byCells(ps)))

	for _, p := range ps {
		colorStyle := lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[p.Color].Cell))

		sb.WriteString(colorStyle.Render("  "))
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%-5d", p.Cells))
		sb.WriteString("  ")
	}
	return sb.String()
}

func (l *Lobby) ViewBoard(top, left, width, height int) string {

	// Arbitrary limits to avoid unreasonable terminal sizes
	// This already shows the board 4 times
	if width > defaultWidth*2 {
		width = defaultHeight * 2
	}
	if height > defaultHeight*2 {
		height = defaultHeight * 2
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
					style = style.Foreground(lipgloss.Color(ColorTable[player.Color].Cursor))

					break
				}
			}

			if l.board[boundY][boundX].Player == life.DeadPlayer && l.board[boundY][boundX].PausedPlayer == life.DeadPlayer && !cursor {
				deadCount++
				continue
			}
			sb.WriteString(deadStyle.Render(strings.Repeat("  ", deadCount)))
			deadCount = 0

			if l.board[boundY][boundX].Player != life.DeadPlayer {
				player, ok := l.players[l.board[boundY][boundX].Player]
				if ok {
					style = style.Background(lipgloss.Color(ColorTable[player.Color].Cell))
				}
			}
			if l.board[boundY][boundX].PausedPlayer != life.DeadPlayer {
				player, ok := l.players[l.board[boundY][boundX].PausedPlayer]
				if ok {
					if !cursor {
						style = style.Foreground(lipgloss.Color(ColorTable[player.Color].Cell))
						pixel = "::"
					} else {
						// style = style.Background(lipgloss.Color(ColorTable[fc].cell))
						pixel = ":]"
					}
				}
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

	if p.Paused {
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
	}

	p.Paused = !p.Paused
}

func (l *Lobby) GetPlayer(id int) *PlayerState {
	// TODO i can't tell if this mutex is rlock is necessary
	l.playersMutex.RLock()
	defer l.playersMutex.RUnlock()
	return l.players[id]
}
