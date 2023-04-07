package middleware

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
	"github.com/zhengkyl/gol/ui/life"
)

type ClientState struct {
	PosX int
	PosY int
	// Pos    Coord
	Paused bool
	Color  int
}
type pixelLookup struct {
	start int
	end   int
}
type Game struct {
	clients      map[*tea.Program]*ClientState
	board        [][]life.Cell
	bufferLookup [][]pixelLookup
	buffer       string
	players      int
	ticker       *time.Ticker
}

const tickRate = 60
const drawRate = 5
const ticksPerDraw = tickRate / drawRate

func NewGame() *Game {

	w, h := 200, 200
	lookup := make([][]pixelLookup, h)
	for y := range lookup {
		lookup[y] = make([]pixelLookup, w)
	}

	return &Game{
		clients:      make(map[*tea.Program]*ClientState),
		board:        life.NewBoard(w, h),
		bufferLookup: lookup,
		players:      0,
		ticker:       time.NewTicker(time.Second / tickRate),
	}
}

func (g *Game) Run() {
	go func() {

		var prevUpdate time.Time
		iteration := 0

		for now := range g.ticker.C {
			iteration++

			if iteration == ticksPerDraw {
				iteration = 0
				g.UpdateBoard()
			}

			g.Update(now.Sub(prevUpdate))

			prevUpdate = now
		}
	}()
}

func (g *Game) Join(p *tea.Program, cs *ClientState) {
	posX := rand.Intn(len(g.board))
	posY := rand.Intn(len(g.board))
	g.players++

	cs.PosX = posX
	cs.PosY = posY
	cs.Paused = true
	cs.Color = g.players
	g.clients[p] = cs

	// TODO testing only!! should only run once
	g.Run()

}

func (g *Game) Leave(p *tea.Program) {
	delete(g.clients, p)
	g.players--
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

var aliveStyle = lipgloss.NewStyle().Background(lipgloss.Color("227"))
var cornerStyle = lipgloss.NewStyle().Background(lipgloss.Color("69"))
var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))

// var aliveWidth = len(aliveStyle.Render(pixel))

func (g *Game) UpdateBoard() {
	g.board = life.NextBoard(g.board)

	// var clients []*ClientState
	// for _, cs := range g.clients {
	// 	clients = append(clients, cs)
	// }
	// sort.Sort(byPos(clients))

	sb := strings.Builder{}

	bytePos := 0

	for y, row := range g.board {
		// line := ""
		for x, cell := range row {
			style := deadStyle
			if cell.IsAlive() {
				style = aliveStyle
			} else if y == 0 && x == 0 {
				style = cornerStyle
			} else if y == len(g.board)-1 && x == len(row)-1 {
				style = cornerStyle
			}

			var pixel = "  "

			for _, cs := range g.clients {
				if y == cs.PosY && x == cs.PosX {
					pixel = "[]"
				}
			}

			g.bufferLookup[y][x].start = bytePos

			sb.WriteString(style.Render(pixel))
			bytePos := sb.Len()

			g.bufferLookup[y][x].end = bytePos
			// line += style.Render(pixel)
		}
		sb.WriteString("\n")
	}
	g.buffer = sb.String()[:sb.Len()-1]

	for p := range g.clients {
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

		repeats := width / boardWidth

		// boundXStart == boundXEndIncl+1 IMPLIES repeats > 0 (at least 1)
		// ALL other cases are inconclusive, either repeats == 0 or repeats > 0
		wrap := boundXStart > (boundXEndIncl+1) || repeats > 0

		// remove 1 when repeat is discontinuous
		if boundXStart <= (boundXEndIncl + 1) {
			repeats--
		}

		if wrap {
			start := g.bufferLookup[boundY][boundXStart].start
			end := g.bufferLookup[boundY][boardWidth-1].end
			sb.WriteString(g.buffer[start:end])

			for i := 0; i < repeats; i++ {
				start := g.bufferLookup[boundY][0].start
				end := g.bufferLookup[boundY][boardWidth-1].end
				sb.WriteString(g.buffer[start:end])
			}

			start = g.bufferLookup[boundY][0].start
			end = g.bufferLookup[boundY][boundXEndIncl].start
			sb.WriteString(g.buffer[start:end])

		} else {
			start := g.bufferLookup[boundY][boundXStart].start
			end := g.bufferLookup[boundY][boundXEndIncl].end
			sb.WriteString(g.buffer[start:end])
		}

		sb.WriteString("\n")
	}

	return sb.String()[:sb.Len()-1]
}

func (g *Game) Update(delta time.Duration) {

}

// https://github.com/charmbracelet/wish/blob/main/bubbletea/tea.go
func GameMiddleware(bth bm.ProgramHandler, cp termenv.Profile) wish.Middleware {
	return func(sh ssh.Handler) ssh.Handler {
		lipgloss.SetColorProfile(cp)

		return func(s ssh.Session) {
			p := bth(s)

			if p != nil {
				_, windowChanges, _ := s.Pty()

				go func() {
					for {
						select {
						case <-s.Context().Done():
							if p != nil {
								p.Quit()
								return
							}
						case w := <-windowChanges:
							if p != nil {
								p.Send(tea.WindowSizeMsg{Width: w.Width, Height: w.Height})
							}
						}
					}
				}()

				if _, err := p.Run(); err != nil {
					log.Error("client exit with error", "error", err)
				}
			}
			sh(s)
		}
	}
}
