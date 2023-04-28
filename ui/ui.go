package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gol/game"
	"github.com/zhengkyl/gol/ui/common"
	"github.com/zhengkyl/gol/ui/keybinds"
	"github.com/zhengkyl/gol/ui/menu"
	"github.com/zhengkyl/gol/ui/multiplayer"
	"github.com/zhengkyl/gol/ui/singleplayer"
)

type screen int

const (
	loadingScreen screen = iota
	menuScreen
	singleplayerScreen
	multiplayerScreen
)

type model struct {
	common   common.Common
	playerId int
	gm       *game.Manager
	menu     *menu.Model
	game     tea.Model
	screen   screen
}

func New(width, height int, gm *game.Manager) model {
	return model{
		screen: loadingScreen,
		common: common.Common{Width: width, Height: height},
		gm:     gm,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

type PlayerId int

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case PlayerId:
		m.playerId = int(msg)

		m.menu = menu.New(m.common, m.gm, m.playerId)
		m.screen = menuScreen

		return m, m.menu.Init()

	case game.JoinSuccessMsg:
		// switch to game view
		m.game = multiplayer.New(common.Common{
			Width: m.common.Width, Height: m.common.Height,
		}, msg)
		m.screen = multiplayerScreen
	case game.SoloGameMsg:

		m.game = singleplayer.New(m.common.Width/2, m.common.Height)

		m.screen = singleplayerScreen
	// switch to solo game view
	case tea.KeyMsg:
		// TODO disconnect or let it handle itself?
		if key.Matches(msg, keybinds.KeyBinds.Quit) {
			// TODO disconnect or let it handle itself?
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.screen {
	case singleplayerScreen:
		_, cmd = m.game.Update(msg)
	case multiplayerScreen:
		_, cmd = m.game.Update(msg)
	case menuScreen:
		_, cmd = m.menu.Update(msg)
	}
	return m, cmd
}

func (m *model) View() string {
	switch m.screen {
	case singleplayerScreen:
		return m.game.View()
	case multiplayerScreen:
		return m.game.View()
	case menuScreen:
		return m.menu.View()
	default:
		return "Loading..."
	}
}
