package ui

import (
	"gol/life"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	boardWidth  int
	boardHeight int
	board       [][]bool
	posX        int
	posY        int
	paused      bool
}

func New(width, height int) model {
	return model{
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return New(msg.Width/2, msg.Height-1), nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "w":
			m.posY = (m.posY - 1 + m.boardHeight) % m.boardHeight
		case "a":
			m.posX = (m.posX - 1 + m.boardWidth) % m.boardWidth
		case "s":
			m.posY = (m.posY + 1 + m.boardHeight) % m.boardHeight
		case "d":
			m.posX = (m.posX + 1 + m.boardWidth) % m.boardWidth
		case " ":
			m.board[m.posY][m.posX] = !m.board[m.posY][m.posX]
		case "enter":
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

func (m model) View() string {

	sb := strings.Builder{}

	for y := range m.board {
		for x, alive := range m.board[y] {

			pixel := "  "
			if y == m.posY && x == m.posX {
				pixel = "[]"
			}

			style := deadStyle
			if alive {
				style = aliveStyle
			}

			sb.WriteString(style.Render(pixel))
		}

		sb.WriteString("\n")
	}

	status := "Playing"
	if m.paused {
		status = "Paused"
	}
	sb.WriteString(status)

	return sb.String()
}
