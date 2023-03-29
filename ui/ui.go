package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gtg/ui/life"
)

type Model struct {
	boardWidth  int
	boardHeight int
	board       [][]life.Life
	posX        int
	posY        int
	paused      bool
}

var aliveStyle = lipgloss.NewStyle().Background(lipgloss.Color("201"))
var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))

func New(width, height int) Model {
	boardWidth := (width / 2)

	// space for mode
	boardHeight := height - 1

	return Model{
		boardWidth:  boardWidth,
		boardHeight: boardHeight,
		board:       life.NewBoard(boardWidth, boardHeight),
		posX:        boardWidth / 2,
		posY:        boardHeight / 2,
		paused:      true,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.paused = true

		m.boardWidth = msg.Width / 2
		m.boardHeight = msg.Height - 1

		m.board = life.NewBoard(m.boardWidth, m.boardHeight)

	case TickMsg:
		if !m.paused {
			cmds = append(cmds, tickOnce())

			m.board = life.NextBoard(m.board)
		}

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
				cmds = append(cmds, tickOnce())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var lines []string

	for y := range m.board {
		line := ""
		for x, alive := range m.board[y] {
			style := deadStyle
			if alive {
				style = aliveStyle
			}

			pixel := "  "
			if m.paused && y == m.posY && x == m.posX {
				pixel = "[]"
			}

			line += style.Render(pixel)
		}

		lines = append(lines, line)
	}

	mode := "Playing"
	if m.paused {
		mode = "Paused"
	}
	lines = append(lines, mode)

	return strings.Join(lines, "\n")
}

type TickMsg struct{}

func tickOnce() tea.Cmd {

	return tea.Tick(time.Second/5, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}
