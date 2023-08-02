package game

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/game/life"
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
	playerCount  int
	board        [][]life.Cell
	boardMutex   sync.RWMutex
	ticker       *time.Ticker
	name         string
	id           int
}

const MaxPlayers = 10
const MaxPlacedCells = 50
const drawRate = 20
const generationRate = 5
const drawsPerGeneration = drawRate / generationRate

const defaultWidth = 160
const defaultHeight = 90

func (l *Lobby) PlayerCount() int {
	// TODO mutex or atomic
	return l.playerCount
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

func (l *Lobby) Join(playerId int, p *tea.Program) (*PlayerState, error) {
	l.playersMutex.Lock()
	defer l.playersMutex.Unlock()

	if l.playerCount == MaxPlayers {
		return nil, fmt.Errorf("Lobby has reached capacity of %v", MaxPlayers)
	}

	l.playerCount++

	posX := rand.Intn(len(l.board))
	posY := rand.Intn(len(l.board))

	var color int
	for i := 1; i <= 11; i++ {
		if !l.playerColors[i] {
			l.playerColors[i] = true
			color = i
			break
		}
	}

	ps := &PlayerState{
		Id:      playerId,
		Program: p,
		PosX:    posX,
		PosY:    posY,
		Paused:  true,
		Color:   color,
	}

	l.players[playerId] = ps

	return ps, nil
}

func (l *Lobby) Leave(playerId int) {

	l.playersMutex.Lock()
	defer l.playersMutex.Unlock()
	l.boardMutex.Lock()
	defer l.boardMutex.Unlock()

	l.playerCount--

	l.playerColors[l.players[playerId].Color] = false
	delete(l.players, playerId)

	for y, row := range l.board {
		for x, cell := range row {
			if cell.Player == playerId {
				l.board[y][x].Player = life.DeadPlayer
				l.board[y][x].PausedPlayer = life.DeadPlayer
			}
		}
	}
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
