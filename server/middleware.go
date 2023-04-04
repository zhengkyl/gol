package server

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
	"github.com/zhengkyl/gol/ui"
	"github.com/zhengkyl/gol/ui/life"
)

type Game struct {
	clients map[*tea.Program]struct{} // set of clients
	board   [][]life.Cell
}

func NewGame() Game {
	return Game{
		clients: make(map[*tea.Program]struct{}),
		// board:
	}
}

func (g Game) Join(p *tea.Program) {
	g.clients[p] = struct{}{}

}
func (g Game) Leave(p *tea.Program) {
	delete(g.clients, p)
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

func teaHandler(s ssh.Session) *tea.Program {
	pty, _, active := s.Pty()

	if !active {
		wish.Fatalln(s, "l + ratio, no active terminal")
		return nil
	}

	ui := ui.New(pty.Window.Width, pty.Window.Height)

	p := tea.NewProgram(ui, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen())

	return p
}

// TODO list
// pass function/pointer/closure to teaprogram
