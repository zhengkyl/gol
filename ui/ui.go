package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gol/server/game"
	"github.com/zhengkyl/gol/ui/keybinds"
)

type model struct {
	clientState    *game.ClientState
	game           *game.Game
	boardWidth     int
	boardHeight    int
	viewportWidth  int
	viewportHeight int
	// board       [][]life.Cell
	// clientState.PosX        int
	// clientState.PosY        int
	// paused      bool
}

// var aliveStyle = lipgloss.NewStyle().Background(lipgloss.Color("227"))
// var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))

type RenderMsg struct{}

type TickMsg struct{}

func tickOnce() tea.Cmd {

	return tea.Tick(time.Second/5, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func New(width, height int, cs *game.ClientState, g *game.Game) model {
	boardWidth, boardHeight := g.BoardSize()

	return model{
		clientState:    cs,
		game:           g,
		boardWidth:     boardWidth,
		boardHeight:    boardHeight,
		viewportWidth:  width / 2,
		viewportHeight: height - 1,
		// board:       life.NewBoard(boardWidth, boardHeight),
		// clientState.PosX:        boardWidth / 2,
		// clientState.PosY:        boardHeight / 2,
		// paused:      true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func mod(dividend, divisor int) int {
	return (dividend + divisor) % divisor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m.viewportWidth = msg.Width / 2
		m.viewportHeight = msg.Height - 1

		m.clientState.Paused = true

	// case TickMsg:
	// 	if !m.clientState.Paused {
	// 		// cmds = append(cmds, tickOnce())

	// 		// m.board = life.NextBoard(m.board)
	// 	}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.KeyBinds.Quit):
			return m, tea.Quit
		case key.Matches(msg, keybinds.KeyBinds.Up):
			m.clientState.PosY = mod(m.clientState.PosY-1, m.boardHeight)
		case key.Matches(msg, keybinds.KeyBinds.Left):
			m.clientState.PosX = (m.clientState.PosX - 1 + m.boardWidth) % m.boardWidth
		case key.Matches(msg, keybinds.KeyBinds.Down):
			m.clientState.PosY = (m.clientState.PosY + 1 + m.boardHeight) % m.boardHeight
		case key.Matches(msg, keybinds.KeyBinds.Right):
			m.clientState.PosX = (m.clientState.PosX + 1 + m.boardWidth) % m.boardWidth
		case key.Matches(msg, keybinds.KeyBinds.Place):
			m.game.Place(m.clientState)
		case key.Matches(msg, keybinds.KeyBinds.Pause):
			m.clientState.Paused = !m.clientState.Paused
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	sb := strings.Builder{}

	sb.WriteString(m.game.ViewBoard(0, 0, m.viewportWidth, m.viewportHeight))
	// sb.WriteString(m.game.ViewBoard(0, 0, 21, 21))
	sb.WriteString("\n")

	help := "wasd/move - <space>/place - <enter>/pause"

	mode := "Playing    "
	if m.clientState.Paused {
		mode = "Paused     "
	}

	mode += help
	sb.WriteString(help)

	return sb.String()
}
