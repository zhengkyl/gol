package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/ui/keybinds"
	"github.com/zhengkyl/gol/ui/life"
)

type model struct {
	boardWidth  int
	boardHeight int
	board       [][]life.Life
	posX        int
	posY        int
	paused      bool
}

var aliveStyle = lipgloss.NewStyle().Background(lipgloss.Color("227"))
var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))

type TickMsg struct{}

func tickOnce() tea.Cmd {

	return tea.Tick(time.Second/5, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func New(width, height int) model {
	boardWidth := (width / 2)

	// space for mode
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.board[m.posY][m.posX] = !m.board[m.posY][m.posX]
		case key.Matches(msg, keybinds.KeyBinds.Pause):
			m.paused = !m.paused
			if !m.paused {
				cmds = append(cmds, tickOnce())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
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

	help := "wasd/move - <space>/place - <enter>/pause"

	mode := "Playing    "
	if m.paused {
		mode = "Paused     "
	}

	mode += help
	lines = append(lines, mode)

	return strings.Join(lines, "\n")
}
