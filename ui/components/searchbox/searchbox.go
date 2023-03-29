package searchbox

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gtg/ui/common"
	"github.com/zhengkyl/gtg/ui/state"
)

type Model struct {
	common common.Common
	state  *state.State
}

func New(c common.Common, s *state.State) *Model {
	return &Model{
		common: c,
		state:  s,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return "ui model"
}
