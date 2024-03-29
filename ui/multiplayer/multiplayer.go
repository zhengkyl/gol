package multiplayer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/game"
	"github.com/zhengkyl/gol/ui/common"
	"github.com/zhengkyl/gol/ui/keybinds"
	"github.com/zhengkyl/gol/util"
)

type model struct {
	playerState *game.PlayerState
	lobby       *game.Lobby
	boardWidth  int
	boardHeight int
	//
	viewportWidth  int
	viewportHeight int
	viewportPosY   int
	viewportPosX   int
}

func New(c common.Common, msg game.JoinSuccessMsg) *model {

	vw := c.Width / 2
	vh := c.Height - 2

	return &model{
		viewportWidth:  vw,
		viewportHeight: vh,

		lobby:        msg.Lobby,
		playerState:  msg.PlayerState,
		boardWidth:   msg.BoardWidth,
		boardHeight:  msg.BoardHeight,
		viewportPosY: util.Mod(msg.PlayerState.PosY-vh/2, msg.BoardHeight),
		viewportPosX: util.Mod(msg.PlayerState.PosX-vw/2, msg.BoardWidth),
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width / 2
		m.viewportHeight = msg.Height - 2

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
			m.playerState.PosY = util.Mod(m.playerState.PosY-1, m.boardHeight)
			if m.playerState.PosY == util.Mod(m.viewportPosY-1, m.boardHeight) {
				m.viewportPosY = m.playerState.PosY
			}

		case key.Matches(msg, keybinds.KeyBinds.Left):
			m.playerState.PosX = util.Mod(m.playerState.PosX-1, m.boardWidth)
			if m.playerState.PosX == util.Mod(m.viewportPosX-1, m.boardWidth) {
				m.viewportPosX = m.playerState.PosX
			}
		case key.Matches(msg, keybinds.KeyBinds.Down):

			m.playerState.PosY = util.Mod(m.playerState.PosY+1, m.boardHeight)
			if m.playerState.PosY == util.Mod(m.viewportPosY+m.viewportHeight, m.boardHeight) {
				m.viewportPosY = util.Mod(m.viewportPosY+1, m.boardHeight)
			}
		case key.Matches(msg, keybinds.KeyBinds.Right):
			m.playerState.PosX = util.Mod(m.playerState.PosX+1, m.boardWidth)
			if m.playerState.PosX == util.Mod(m.viewportPosX+m.viewportWidth, m.boardWidth) {
				m.viewportPosX = util.Mod(m.viewportPosX+1, m.boardWidth)
			}

		case key.Matches(msg, keybinds.KeyBinds.Place):
			m.lobby.Place(m.playerState.Id)
		case key.Matches(msg, keybinds.KeyBinds.Enter):
			m.lobby.TogglePause(m.playerState.Id)
		}
	}

	return m, nil
}

var (
	helpStyle = lipgloss.NewStyle().Inline(true)
)

func (m *model) View() string {
	if m.lobby == nil {
		return "loading... probably a critical error"
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

	sb.WriteString(helpStyle.MaxWidth(m.viewportWidth*2).Render(
		"wasd/hjkl/←↑↓→",
		"move",
		" • ",
		"<space>",
		"place",
		" • ",
		"<enter>",
		"play/edit",
		" • ",
		"<esc>",
		"menu",
	))

	return sb.String()
}
