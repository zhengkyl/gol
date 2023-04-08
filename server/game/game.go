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
type pixelLookup struct {
	start int
	end   int
}

type GameState int

const (
	PAUSED GameState = iota
	PLAYING
)

type Game struct {
	clients           PlayerMap
	board             [][]life.Cell
	boardBufferLookup [][]pixelLookup
	boardBuffer       string
	// players           int
	paused int
	ticker *time.Ticker
	state  GameState
}

const MaxPlayers = 10
const MaxPlacedCells = 20
const drawRate = 10
const generationRate = 5
const drawsPerGeneration = drawRate / generationRate

func (g *Game) Players() int {
	return g.clients.Len()
}

func (g *Game) PausedPlayers() int {
	return g.paused
}

func NewGame() *Game {

	w, h := 20, 20
	lookup := make([][]pixelLookup, h)
	for y := range lookup {
		lookup[y] = make([]pixelLookup, w)
	}

	return &Game{
		clients:           PlayerMap{v: make(map[*tea.Program]*ClientState)},
		board:             life.NewBoard(w, h),
		boardBufferLookup: lookup,
		// players:           0,
		paused: 0,
		ticker: time.NewTicker(time.Second / drawRate),
		state:  PAUSED,
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

	sb := strings.Builder{}

	for y, row := range g.board {
		// line := ""
		for x, cell := range row {
			style := deadStyle
			if cell.IsAlive() {
				style = lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[cell.Color].cell))
				// store user id so can access them?
			}

			var pixel = "  "
			if y == 0 && x == 0 {
				// Keep track of tiling
				pixel = "::"
			}

			for _, cs := range css {
				if y == cs.PosY && x == cs.PosX {
					pixel = "[]"
					style = style.Copy().Foreground(lipgloss.Color(ColorTable[cs.Color].cursor))
				}
			}

			g.boardBufferLookup[y][x].start = sb.Len()
			sb.WriteString(style.Render(pixel))
			g.boardBufferLookup[y][x].end = sb.Len()
		}
		sb.WriteString("\n")
	}
	g.boardBuffer = sb.String()[:sb.Len()-1]

	for _, p := range ps {
		p.Send(ServerRedrawMsg{})
	}
}

func (g *Game) ViewBoard(top, left, width, height int) string {
	sb := strings.Builder{}

	boardWidth, boardHeight := g.BoardSize()

	for y := top; y < top+height; y++ {

		boundY := (y + boardHeight) % boardHeight

		boundXStart := (left + boardWidth) % boardWidth
		boundXEndIncl := (left + width - 1 + boardWidth) % boardWidth

		// if width % boardWidth == 0 -> special case
		repeats := (width - 1) / boardWidth

		wrap := boundXStart > boundXEndIncl || repeats > 0

		// remove 1 when repeat is discontinuous
		if boundXStart <= boundXEndIncl && repeats > 0 {
			repeats--
		}

		if wrap {

			start := g.boardBufferLookup[boundY][boundXStart].start
			end := g.boardBufferLookup[boundY][boardWidth-1].end
			sb.WriteString(g.boardBuffer[start:end])

			if repeats > 0 {
				start = g.boardBufferLookup[boundY][0].start
				end = g.boardBufferLookup[boundY][boardWidth-1].end
				sb.WriteString(strings.Repeat(g.boardBuffer[start:end], repeats))
			}

			start = g.boardBufferLookup[boundY][0].start
			end = g.boardBufferLookup[boundY][boundXEndIncl].end
			sb.WriteString(g.boardBuffer[start:end])

		} else {
			start := g.boardBufferLookup[boundY][boundXStart].start
			end := g.boardBufferLookup[boundY][boundXEndIncl].end
			sb.WriteString(g.boardBuffer[start:end])
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
