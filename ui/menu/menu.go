package menu

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gol/server/game"
)

type model struct {
	gm *game.Manager
}

func New(gm *game.Manager) *model {
	return &model{gm}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *model) View() string {
	sb := strings.Builder{}

	// gm.

	return ""
}
