package common

import tea "github.com/charmbracelet/bubbletea"

type Common struct {
	Width  int
	Height int
}

type CommonModel interface {
	tea.Model
	SetSize(w, h int)
}

// type DefaultModel struct {
// 	Common Common
// }

// func (m *DefaultModel) SetSize(w, h int) {
// 	m.Common.Width = w
// 	m.Common.Height = h
// }

// func (m *DefaultModel) Init() tea.Cmd {
// 	return nil
// }

// func (m *DefaultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	return m, nil
// }

// func (m *DefaultModel) View() string {
// 	return "ui model"
// }
