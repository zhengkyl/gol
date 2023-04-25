package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gol/game"
	"github.com/zhengkyl/gol/ui/common"
	"github.com/zhengkyl/gol/ui/keybinds"
	"github.com/zhengkyl/gol/ui/menu"
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
	manager  game.Manager
	menu     *menu.Model
	game     tea.Model
	screen   screen
}

func New(width, height int) model {
	return model{
		screen: loadingScreen,
		common: common.Common{Width: width, Height: height},
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

		gm := game.NewManager()
		m.menu = menu.New(m.common, gm, m.playerId)

	case game.JoinSuccessMsg:
		// switch to game view
	// case game.SoloGameMsg:
	// switch to solo game view
	case tea.KeyMsg:
		if key.Matches(msg, keybinds.KeyBinds.Quit) {
			// TODO disconnect or let it handle itself?
			return m, tea.Quit
		}
	}

	switch m.screen {
	case singleplayerScreen:
		return m.game.Update(msg)
	case multiplayerScreen:
		return m.game.Update(msg)
	case menuScreen:
		return m.menu.Update(msg)
	default:
		return m, nil
	}
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
