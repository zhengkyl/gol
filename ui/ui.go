package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/server/game"
	"github.com/zhengkyl/gol/ui/keybinds"
)

type model struct {
	playerState *game.PlayerState
	lobby       *game.Lobby
	id          int
	boardWidth  int
	boardHeight int
	//
	viewportWidth  int
	viewportHeight int
	viewportPosY   int
	viewportPosX   int
	//
}

func New(width, height int) model {

	vw := width / 2
	vh := height - 2

	return model{
		viewportWidth:  vw,
		viewportHeight: vh,
		// asdf: help.New()
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
		m.lobby = msg.Lobby
		m.id = msg.Id
		m.playerState = msg.PlayerState
		m.boardWidth = msg.BoardWidth
		m.boardHeight = msg.BoardHeight

		m.viewportPosY = mod(m.playerState.PosY-m.viewportHeight/2, m.boardHeight)
		m.viewportPosX = mod(m.playerState.PosX-m.viewportWidth/2, m.boardWidth)
		return m, nil

	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width / 2
		m.viewportHeight = msg.Height - 2
		// if m.playerState != nil {
		// 	m.viewportPosY = mod(m.playerState.PosY-m.viewportHeight/2, m.boardHeight)
		// 	m.viewportPosX = mod(m.playerState.PosX+m.viewportWidth/2, m.boardWidth)
		// }

	case game.UpdateBoardMsg:
		// do nothing? just rerender
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.KeyBinds.Quit):
			return m, tea.Quit
		}

		if m.lobby == nil {
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
			m.lobby.Place(m.id)
		case key.Matches(msg, keybinds.KeyBinds.Pause):
			m.lobby.TogglePause(m.id)
		}
	}

	return m, tea.Batch(cmds...)
}

var (
	helpStyle = lipgloss.NewStyle().Inline(true)
)

func (m *model) View() string {
	if m.lobby == nil {
		return "loading... probably a critical error, there should be negligible loading"
	}

	sb := strings.Builder{}

	avatarStyle := lipgloss.NewStyle().Background(lipgloss.Color(game.ColorTable[m.playerState.Color].Cell))

	mode := "PLAYING"
	if m.playerState.Paused {
		mode = fmt.Sprintf("EDITING %d/%d cells placed", m.playerState.Placed, game.MaxPlacedCells)
	}

	sb.WriteString(helpStyle.MaxWidth(m.viewportWidth*2).Render(
		avatarStyle.Render("  "),
		fmt.Sprintf("%-30s", mode),
		"SCORE",
		m.lobby.Scoreboard(),
	))

	sb.WriteString("\n")
	sb.WriteString(m.lobby.ViewBoard(m.viewportPosY, m.viewportPosX, m.viewportWidth, m.viewportHeight))
	sb.WriteString("\n")

	// helpSb stri

	sb.WriteString(helpStyle.MaxWidth(m.viewportWidth*2).Render(
		"wasd/hjkl/←↑↓→",
		"move",
		" • ",
		"<space>",
		"place",
		" • ",
		"<enter>",
		"play/edit",
	))

	// mode += help
	// mode += fmt.Sprintf("            %d/%d cells placed", m.playerState.Placed, game.MaxPlacedCells)
	// mode += fmt.Sprintf("            %d/%d players", m.lobby.PlayerCount(), game.MaxPlayers)
	// mode += fmt.Sprintf("            %s", m.playerState.Test)

	// sb.WriteString(mode)

	return sb.String()
}
