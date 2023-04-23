package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gol/game"
	"github.com/zhengkyl/gol/server"
	"github.com/zhengkyl/gol/ui/common"
	"github.com/zhengkyl/gol/ui/keybinds"
	"github.com/zhengkyl/gol/ui/menu"
)

type screen int

const (
	menuScreen screen = iota
	singleplayerScreen
	multiplayerScreen
)

type model struct {
	playerId int
	manager  game.Manager
	menu     *menu.Model
	game     tea.Model
	screen   screen
}

func New(width, height int) model {
	common := common.Common{
		Width:  width,
		Height: height,
	}
	gm := game.NewManager()

	return model{
		menu:   menu.New(common, gm),
		screen: menuScreen,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case server.PlayerId:
		m.playerId = int(msg)
	case game.JoinLobbyMsg:
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
	default:
		return m.menu.Update(msg)
	}
}

func (m *model) View() string {
	switch m.screen {
	case singleplayerScreen:
		return m.game.View()
	case multiplayerScreen:
		return m.game.View()
	default:
		return m.menu.View()
	}
}
