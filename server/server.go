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
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
	"github.com/zhengkyl/gol/server/middleware"
	"github.com/zhengkyl/gol/ui"
)

const (
	host = "0.0.0.0"
	port = 2345
)

func RunServer() {

	game := middleware.NewGame()

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/server_ed25519"),
		wish.WithMiddleware(
			middleware.GameMiddleware(teaHandler(game), termenv.ANSI256),
			// bm.Middleware(teaHandler),
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

func teaHandler(game *middleware.Game) bm.ProgramHandler {
	return func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()

		if !active {
			wish.Fatalln(s, "l + ratio, no active terminal")
			return nil
		}

		cs := &middleware.ClientState{}

		ui := ui.New(pty.Window.Width, pty.Window.Height, cs, game)

		p := tea.NewProgram(ui, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen())

		game.Join(p, cs)

		return p
	}
}
