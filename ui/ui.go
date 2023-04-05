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
	boardWidth := width / 2
	boardHeight := height - 1

	return model{
		boardWidth:  boardWidth,
		boardHeight: boardHeight,
		board:       life.NewBoard(boardWidth, boardHeight),
		posX:        boardWidth / 2,
		posY:        boardHeight / 2,
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
		return New(msg.Width, msg.Height), nil

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

	var lines []string

	for y := range m.board {
		line := ""
		for x, alive := range m.board[y] {

			pixel := "  "
			if y == m.posY && x == m.posX {
				pixel = "[]"
			}

			style := deadStyle
			if alive {
				style = aliveStyle
			}

			line += style.Render(pixel)
		}

		lines = append(lines, line)
	}

	status := "Playing"
	if m.paused {
		status = "Paused"
	}
	lines = append(lines, status)

	return strings.Join(lines, "\n")
}
