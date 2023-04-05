package middleware

import (
	"math/rand"
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
	posX   int
	posY   int
	paused bool
	color  int
}

type Game struct {
	clients map[*tea.Program]ClientState
	board   [][]life.Cell
	players int
	ticker  *time.Ticker
}

const tickRate = 60
const drawRate = 5
const ticksPerDraw = tickRate / drawRate

func NewGame() Game {

	return Game{
		clients: make(map[*tea.Program]ClientState),
		board:   life.NewBoard(200, 200),
		players: 0,
		ticker:  time.NewTicker(time.Second / tickRate),
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

func (g *Game) Join(p *tea.Program) {
	posX := rand.Intn(len(g.board))
	posY := rand.Intn(len(g.board))
	g.players++

	g.clients[p] = ClientState{
		posX:   posX,
		posY:   posY,
		paused: true,
		color:  g.players,
	}

}

func (g *Game) Leave(p *tea.Program) {
	delete(g.clients, p)
	g.players--
}

func (g *Game) Move() {
	// update
}

func (g *Game) Unpause() {
	// update
}

func (g *Game) Iterate() {
	// update
}

type ServerRedrawMsg struct{}

func (g *Game) UpdateBoard() {
	g.board = life.NextBoard(g.board)

	for p, cs := range g.clients {
		p.Send(ServerRedrawMsg{})
	}
}

func (g *Game) Update(delta time.Duration) {

}

// https://github.com/charmbracelet/wish/blob/main/bubbletea/tea.go
func gameMiddleware(bth bm.ProgramHandler, cp termenv.Profile) wish.Middleware {
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

// TODO list
// pass function/pointer/closure to teaprogram
