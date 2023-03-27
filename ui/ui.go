package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gtg/ui/common"
)

type Model struct {
	common common.Common
}

func New(width, height int) *Model {
	return &Model{
		common.Common{
			Width:  width,
			Height: height,
		},
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
