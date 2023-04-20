package menu

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/game"
	"github.com/zhengkyl/gol/ui/common"
)

type model struct {
	gm      *game.Manager
	common  common.Common
	options []listItem
}

func New(gm *game.Manager) *model {
	return &model{gm: gm}
}

func (m *model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []game.LobbyStatus:
		for _, status := range msg {

		}
	}
	return m, nil
}

func combine(left, right string, width int) string {
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)

	spaces := width - (leftW + rightW)

	if spaces < 1 {
		if leftW > width {
			return left[:width-1] + "…"
		}
		return left
	} else {
		return left + strings.Repeat(" ", spaces) + right
	}
}

func (m *model) View() string {
	sb := strings.Builder{}

	for _, li := range m.options {
		sb.WriteString(combine(li.titleLeft, li.titleRight, m.common.Width))
		sb.WriteString("\n")

		sb.WriteString(combine(li.descLeft, li.descRight, m.common.Width))
	}
	// gm.

	return ""
}

type listItem struct {
	titleLeft  string
	titleRight string
	descLeft   string
	descRight  string
}
