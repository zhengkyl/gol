package server

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/zhengkyl/gtg/ui"
)

const (
	host = "localhost"
	port = 2634
)

func NewServer() *ssh.Server {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
		),
	)

	if err != nil {
		log.Error("server didn't start", "err", err)
	}

	return s
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()

	if !active {
		wish.Fatalln(s, "l + ratio, no active terminal")
		return nil, nil
	}

	ui := ui.New(pty.Window.Width, pty.Window.Height)

	return ui, []tea.ProgramOption{tea.WithAltScreen()}
}
