package singleplayer

import (
	"strings"
	"time"

	"github.com/zhengkyl/gol/game/life"
	"github.com/zhengkyl/gol/ui/keybinds"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	dead   = 0
	player = 1
)

type model struct {
	boardWidth  int
	boardHeight int
	board       [][]life.Cell
	posX        int
	posY        int
	paused      bool
}

func New(width, height int) *model {
	return &model{
		boardWidth:  width,
		boardHeight: height,
		board:       life.NewBoard(width, height),
		posX:        width / 2,
		posY:        height / 2,
		paused:      true,
	}
}

type tickMsg struct{}

var tickOnce = tea.Tick(time.Second/5, func(t time.Time) tea.Msg {
	return tickMsg{}
})

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return New(msg.Width/2, msg.Height-1), nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.KeyBinds.Quit):
			return m, tea.Quit
		case key.Matches(msg, keybinds.KeyBinds.Up):
			m.posY = (m.posY - 1 + m.boardHeight) % m.boardHeight
		case key.Matches(msg, keybinds.KeyBinds.Left):
			m.posX = (m.posX - 1 + m.boardWidth) % m.boardWidth
		case key.Matches(msg, keybinds.KeyBinds.Down):
			m.posY = (m.posY + 1 + m.boardHeight) % m.boardHeight
		case key.Matches(msg, keybinds.KeyBinds.Right):
			m.posX = (m.posX + 1 + m.boardWidth) % m.boardWidth
		case key.Matches(msg, keybinds.KeyBinds.Place):
			if m.board[m.posY][m.posX].Player == dead {
				m.board[m.posY][m.posX].Player = player
			} else {
				m.board[m.posY][m.posX].Player = dead
			}
		case key.Matches(msg, keybinds.KeyBinds.Enter):
			m.paused = !m.paused
			if !m.paused {
				return m, tickOnce
			}
		}

	case tickMsg:
		if !m.paused {
			m.board = life.NextBoard(m.board)
			return m, tickOnce
		}
	}

	return m, nil
}

var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))
var aliveStyle = lipgloss.NewStyle().Background(lipgloss.Color("227"))

func (m *model) View() string {

	sb := strings.Builder{}

	for y := range m.board {
		for x, cell := range m.board[y] {

			pixel := "  "
			if y == m.posY && x == m.posX {
				pixel = "[]"
			}

			style := deadStyle
			if cell.Player == player {
				style = aliveStyle
			}

			sb.WriteString(style.Render(pixel))
		}

		sb.WriteString("\n")
	}

	status := "Playing"
	if m.paused {
		status = "Paused "
	}
	sb.WriteString(status + "  •  wasd/hjkl/←↑↓→ move  •  <space> place  •  <enter> play/pause  •  <esc> menu")
	return sb.String()
}
