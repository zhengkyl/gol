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
	playerState    *game.PlayerState
	game           *game.Lobby
	id             int
	boardWidth     int
	boardHeight    int
	viewportWidth  int
	viewportHeight int
	viewportPosY   int
	viewportPosX   int
}

func New(width, height int) model {

	vw := width / 2
	vh := height - 1

	return model{
		viewportWidth:  vw,
		viewportHeight: vh,
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
	case game.JoinLobbyMsg:
		m.game = msg.Lobby
		m.id = msg.Id
		m.playerState = m.game.GetPlayer(m.id)
		m.boardWidth, m.boardHeight = m.game.BoardSize()
		m.viewportPosX = 0
		m.viewportPosY = 0

		m.viewportPosY = mod(m.playerState.PosY-m.viewportHeight/2, m.boardHeight)
		m.viewportPosX = mod(m.playerState.PosX+m.viewportWidth/2, m.boardWidth)
		return m, nil

	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width / 2
		m.viewportHeight = msg.Height - 1

		// m.playerState.Paused = true

	case game.ServerRedrawMsg:

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.KeyBinds.Quit):
			return m, tea.Quit
		}

		if m.game == nil {
			return m, nil
		}

		switch {
		case key.Matches(msg, keybinds.KeyBinds.Up):
			m.playerState.PosY = mod(m.playerState.PosY-1, m.boardHeight)
			if m.playerState.PosY == mod(m.viewportPosY-1, m.boardHeight) {
				m.viewportPosY = m.playerState.PosY
			}

		case key.Matches(msg, keybinds.KeyBinds.Left):
			m.playerState.PosX = mod(m.playerState.PosX-1, m.boardWidth)
			if m.playerState.PosX == mod(m.viewportPosX-1, m.boardWidth) {
				m.viewportPosX = m.playerState.PosX
			}
		case key.Matches(msg, keybinds.KeyBinds.Down):

			m.playerState.PosY = mod(m.playerState.PosY+1, m.boardHeight)
			if m.playerState.PosY == mod(m.viewportPosY+m.viewportHeight, m.boardHeight) {
				m.viewportPosY = mod(m.viewportPosY+1, m.boardHeight)
			}
		case key.Matches(msg, keybinds.KeyBinds.Right):
			m.playerState.PosX = mod(m.playerState.PosX+1, m.boardWidth)
			if m.playerState.PosX == mod(m.viewportPosX+m.viewportWidth, m.boardWidth) {
				m.viewportPosX = mod(m.viewportPosX+1, m.boardWidth)
			}

		case key.Matches(msg, keybinds.KeyBinds.Place):
			m.game.Place(m.id)
		case key.Matches(msg, keybinds.KeyBinds.Pause):
			m.playerState.Paused = !m.playerState.Paused
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.game == nil {
		return "loading"
	}

	sb := strings.Builder{}

	sb.WriteString(m.game.ViewBoard(m.viewportPosY, m.viewportPosX, m.viewportWidth, m.viewportHeight))
	sb.WriteString("\n")

	help := "wasd/move - <space>/place - <enter>/pause"

	mode := "Playing    "
	if m.playerState.Paused {
		mode = "Paused     "
	}

	mode += help
	mode += fmt.Sprintf("            %d/%d cells placed", m.playerState.Placed, game.MaxPlacedCells)
	// mode += fmt.Sprintf("            %s", m.playerState.Test)

	sb.WriteString(mode)

	return sb.String()
}
