package game

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/ui/life"
)

type ClientState struct {
	PosX int
	PosY int
	// Pos    Coord
	Paused bool
	Color  int
	Placed int
	Cells  int
}

type GameState int

const (
	PAUSED GameState = iota
	PLAYING
)

type Game struct {
	clients PlayerMap
	board   [][]life.Cell
	paused  int
	ticker  *time.Ticker
	state   GameState
}

const MaxPlayers = 10
const MaxPlacedCells = 20
const drawRate = 10
const generationRate = 5
const drawsPerGeneration = drawRate / generationRate

const size = 100

func (g *Game) Players() int {
	return g.clients.Len()
}

func (g *Game) PausedPlayers() int {
	return g.paused
}

func NewGame() *Game {
	w, h := size, size

	return &Game{
		clients: PlayerMap{v: make(map[*tea.Program]*ClientState)},
		board:   life.NewBoard(w, h),
		paused:  0,
		ticker:  time.NewTicker(time.Second / drawRate),
		state:   PAUSED,
	}
}

func (g *Game) Run() {
	go func() {

		var prevUpdate time.Time
		iteration := 0

		for now := range g.ticker.C {
			iteration++

			if iteration == drawsPerGeneration {
				iteration = 0
				g.UpdateBoard()
			}

			g.Update(now.Sub(prevUpdate))

			prevUpdate = now
		}
	}()
}

func (g *Game) Join(p *tea.Program, cs *ClientState) bool {
	if g.clients.Len() == MaxPlayers {
		return false
	}

	posX := rand.Intn(len(g.board))
	posY := rand.Intn(len(g.board))

	cs.PosX = posX
	cs.PosY = posY
	cs.Paused = true
	cs.Color = g.clients.Len() + 1
	g.clients.Set(p, cs)

	return true
}

func (g *Game) Leave(p *tea.Program) {
	g.clients.Delete(p)
}

func (g *Game) Unpause() {
	// update
}

func (g *Game) Iterate() {
	// update
}

func (g *Game) BoardSize() (int, int) {
	return len(g.board), len(g.board[0])
}

type ServerRedrawMsg struct{}

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

func (g *Game) UpdateBoard() {

	if g.state != PLAYING {
		return
	}

	g.board = life.NextBoard(g.board)
}

func (g *Game) Update(delta time.Duration) {
	g.paused = 0

	ps, css := g.clients.Entries()

	for _, cs := range css {
		if cs.Paused {
			g.paused++
		}
	}

	if g.paused > len(ps)/2 {
		g.state = PAUSED
	} else {
		g.state = PLAYING
	}

	for _, p := range ps {
		p.Send(ServerRedrawMsg{})
	}
}

func (g *Game) ViewBoard(top, left, width, height int) string {

	// Arbitrary limits to avoid unreasonable terminal sizes
	// This already shows the board 4 times
	if width > size*2 {
		width = size * 2
	}
	if height > size*2 {
		height = size * 2
	}

	sb := strings.Builder{}

	_, css := g.clients.Entries()

	boardWidth, boardHeight := g.BoardSize()

	for y := top; y < top+height; y++ {
		boundY := (y + boardHeight) % boardHeight

		deadCount := 0
		for x := left; x < left+width; x++ {
			boundX := (x + boardWidth) % boardWidth
			style := lipgloss.NewStyle()
			pixel := "  "

			for _, cs := range css {
				if boundY == cs.PosY && boundX == cs.PosX {
					pixel = "[]"
					style = style.Foreground(lipgloss.Color(ColorTable[cs.Color].cursor))
				}
			}

			if !g.board[boundY][boundX].IsAlive() && pixel == "  " {
				deadCount++
				continue
			}
			sb.WriteString(deadStyle.Render(strings.Repeat("  ", deadCount)))
			deadCount = 0

			style = style.Background(lipgloss.Color(ColorTable[g.board[boundY][boundX].Color].cell))
			sb.WriteString(style.Render(pixel))
		}

		if deadCount > 0 {
			sb.WriteString(deadStyle.Render(strings.Repeat("  ", deadCount)))
		}

		sb.WriteString("\n")
	}

	return sb.String()[:sb.Len()-1]
}

func (g *Game) Place(cs *ClientState) {

	if !g.board[cs.PosY][cs.PosX].IsAlive() {
		if cs.Placed >= MaxPlacedCells {
			return
		}
		g.board[cs.PosY][cs.PosX].Color = cs.Color
		cs.Placed++
	} else if g.board[cs.PosY][cs.PosX].Color == cs.Color {
		g.board[cs.PosY][cs.PosX].Color = 0
		cs.Placed--
	}
}
