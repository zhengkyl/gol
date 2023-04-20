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
	"github.com/zhengkyl/gol/game"
	"github.com/zhengkyl/gol/ui/menu"
)

const (
	host = "0.0.0.0"
	port = 2345
)

func RunServer() {

	gm := game.NewManager()

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/server_ed25519"),
		wish.WithMiddleware(
			MiddlewareWithProgramHandler(teaHandler(gm), termenv.ANSI256, gm),
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

func teaHandler(gm *game.Manager) bm.ProgramHandler {
	return func(s ssh.Session) *tea.Program {
		_, _, active := s.Pty()
		// pty, _, active := s.Pty()

		if !active {
			wish.Fatalln(s, "l + ratio, no active terminal")
			return nil
		}

		// ui := ui.New(pty.Window.Width, pty.Window.Height)
		// p := tea.NewProgram(&ui, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen())
		menu := menu.New(gm)
		p := tea.NewProgram(menu, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen())

		gm.Connect(p)

		// go func() {
		// 	l := gm.FindLobby()
		// 	id, ok := l.Join(p)

		// 	if !ok {
		// 		wish.Fatalln(s, fmt.Sprintf("Failed to join. %d/%d players in game. :/", l.PlayerCount(), game.MaxPlayers))
		// 		return
		// 	}

		// 	bw, bh := l.BoardSize()

		// 	ps := l.GetPlayer(id)

		// 	if ps == nil {
		// 		wish.Fatalln(s, "Failed to join. PlayerState not initialized.")
		// 		return
		// 	}

		// 	p.Send(game.JoinLobbyMsg{
		// 		Lobby:       l,
		// 		PlayerState: ps,
		// 		Id:          id,
		// 		BoardWidth:  bw,
		// 		BoardHeight: bh,
		// 	})

		// 	s.Context().SetValue("lobby", l)
		// 	s.Context().SetValue("playerId", id)
		// }()

		return p
	}
}

// copied from wish/bubbletea b/c need to know when p.Quit() in order to trigger game.Leave()
func MiddlewareWithProgramHandler(bth bm.ProgramHandler, cp termenv.Profile, gm *game.Manager) wish.Middleware {
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
								gm.Disconnect(p)
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
