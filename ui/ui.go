package ui

import (
	"fmt"
	"strings"

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
	viewportPosY   int
	viewportPosX   int
}

func New(width, height int, cs *game.ClientState, g *game.Game) model {
	boardWidth, boardHeight := g.BoardSize()

	vw := width / 2
	vh := height - 1

	return model{
		clientState:    cs,
		game:           g,
		boardWidth:     boardWidth,
		boardHeight:    boardHeight,
		viewportWidth:  vw,
		viewportHeight: vh,
		viewportPosY:   mod(cs.PosY-vh/2, boardHeight),
		viewportPosX:   mod(cs.PosX+vw/2, boardWidth),
		// board:       life.NewBoard(boardWidth, boardHeight),
		// clientState.PosX:        boardWidth / 2,
		// clientState.PosY:        boardHeight / 2,
		// paused:      true,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func mod(dividend, divisor int) int {
	return (dividend + divisor) % divisor
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m.viewportWidth = msg.Width / 2
		m.viewportHeight = msg.Height - 1

		m.clientState.Paused = true

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.KeyBinds.Quit):
			return m, tea.Quit

		case key.Matches(msg, keybinds.KeyBinds.Up):

			m.clientState.PosY = mod(m.clientState.PosY-1, m.boardHeight)
			if m.clientState.PosY == mod(m.viewportPosY-1, m.boardHeight) {
				m.viewportPosY = m.clientState.PosY
			}

		case key.Matches(msg, keybinds.KeyBinds.Left):
			m.clientState.PosX = mod(m.clientState.PosX-1, m.boardWidth)
			if m.clientState.PosX == mod(m.viewportPosX-1, m.boardWidth) {
				m.viewportPosX = m.clientState.PosX
			}

		case key.Matches(msg, keybinds.KeyBinds.Down):

			m.clientState.PosY = mod(m.clientState.PosY+1, m.boardHeight)
			if m.clientState.PosY == mod(m.viewportPosY+m.viewportHeight, m.boardHeight) {
				m.viewportPosY = mod(m.viewportPosY+1, m.boardHeight)
			}

		case key.Matches(msg, keybinds.KeyBinds.Right):
			m.clientState.PosX = mod(m.clientState.PosX+1, m.boardWidth)
			if m.clientState.PosX == mod(m.viewportPosX+m.viewportWidth, m.boardWidth) {
				m.viewportPosX = mod(m.viewportPosX+1, m.boardWidth)
			}
		case key.Matches(msg, keybinds.KeyBinds.Place):
			m.game.Place(m.clientState)
		case key.Matches(msg, keybinds.KeyBinds.Pause):
			m.clientState.Paused = !m.clientState.Paused
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	sb := strings.Builder{}

	sb.WriteString(m.game.ViewBoard(m.viewportPosY, m.viewportPosX, m.viewportWidth, m.viewportHeight))
	// sb.WriteString(m.game.ViewBoard(0, 0, 21, 21))
	sb.WriteString("\n")

	help := "wasd/move - <space>/place - <enter>/pause"

	mode := "Playing    "
	if m.clientState.Paused {
		mode = "Paused     "
	}

	mode += help
	mode += fmt.Sprintf("            %d/%d cells placed", m.clientState.Placed, game.MaxPlacedCells)
	mode += fmt.Sprintf("            %d/%d players paused", m.game.PausedPlayers(), m.game.Players())

	sb.WriteString(mode)

	return sb.String()
}
