package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
	"github.com/zhengkyl/gol/server/game"
	"github.com/zhengkyl/gol/ui"
)

const (
	host = "0.0.0.0"
	port = 2345
)

func RunServer() {

	game := game.NewGame()

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/server_ed25519"),
		wish.WithMiddleware(
			MiddlewareWithProgramHandler(teaHandler(game), termenv.ANSI256, game),
			lm.Middleware(),
		),
	)

	if err != nil {
		log.Error("server didn't start", "err", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		game.Run()
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
		}
	}()

	<-done
	log.Info("Stopping SSH server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func teaHandler(g *game.Game) bm.ProgramHandler {
	return func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()

		if !active {
			wish.Fatalln(s, "l + ratio, no active terminal")
			return nil
		}

		cs := &game.ClientState{}

		ui := ui.New(pty.Window.Width, pty.Window.Height, cs, g)

		p := tea.NewProgram(&ui, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen())

		if !g.Join(p, cs) {
			wish.Fatalln(s, fmt.Sprintf("Failed to join. %d/%d players in game. :/", g.Players(), game.MaxPlayers))
			return nil
		}

		return p
	}
}

// copied from wish/bubbletea b/c need to know when p.Quit() in order to trigger game.Leave()
func MiddlewareWithProgramHandler(bth bm.ProgramHandler, cp termenv.Profile, game *game.Game) wish.Middleware {
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
								game.Leave(p)
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
					log.Error("app exit with error", "error", err)
				}
				// p.Kill() will force kill the program if it's still running,
				// and restore the terminal to its original state in case of a
				// tui crash
				p.Kill()
			}
			sh(s)
		}
	}
}
