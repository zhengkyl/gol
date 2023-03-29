package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gtg/ui/common"
)

type Model struct {
	common common.Common
	paused bool
	posX   int
	posY   int
	board  [][]bool
	count  int
}

var aliveStyle = lipgloss.NewStyle().Background(lipgloss.Color("201"))
var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))

func New(width, height int) Model {
	boardWidth := width / 2

	board := make([][]bool, height)
	for i := range board {
		board[i] = make([]bool, boardWidth)
	}

	return Model{
		common: common.Common{
			Width:  width,
			Height: height,
		},
		paused: true,
		board:  board,
		posX:   boardWidth / 2,
		posY:   height / 2,
		count:  0,
	}
}

func (m Model) Init() tea.Cmd {
	// return m.tickOnce()
	return nil
}

var dir = [3]int{-1, 0, 1}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.count++
		m.paused = true
		// TODO
		m.common.Width = msg.Width / 2
		m.common.Height = msg.Height

		m.board = make([][]bool, msg.Height)
		for i := range m.board {
			m.board[i] = make([]bool, msg.Width/2)
		}

	case TickMsg:
		if !m.paused {
			cmds = append(cmds, m.tickOnce())

			newBoard := make([][]bool, m.common.Height)
			for i := range newBoard {
				newBoard[i] = make([]bool, m.common.Width)
			}

			for y := range m.board {

				for x := range m.board[y] {

					neighbors := 0

					for _, dirX := range dir {
						if x+dirX < 0 || x+dirX >= len(m.board[y]) {
							continue
						}
						for _, dirY := range dir {
							if y+dirY < 0 || y+dirY >= len(m.board) {
								continue
							}
							if dirY == 0 && dirX == 0 {
								continue
							}

							if m.board[y+dirY][x+dirX] {
								neighbors++
							}

						}
					}

					if !m.board[y][x] && neighbors == 3 {
						newBoard[y][x] = true
					}
					if m.board[y][x] {
						if neighbors < 2 || neighbors > 3 {
							newBoard[y][x] = false
						} else {
							newBoard[y][x] = true
						}
					}

				}
			}

			m.board = newBoard

		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.paused = !m.paused
			if !m.paused {
				cmds = append(cmds, m.tickOnce())
			}
		case " ":
			m.board[m.posY][m.posX] = !m.board[m.posY][m.posX]
		case "w":
			m.posY--
		case "a":
			m.posX--
		case "s":
			m.posY++
		case "d":
			m.posX++
		}
		// cmds = append(cmds, func() tea.Msg { return nil })
		m.count++
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// view := fmt.Sprint(m.count)
	// view := ""
	var lines []string

	for y := range m.board {
		// view += fmt.Sprint(i)
		line := ""
		for x, alive := range m.board[y] {
			var style lipgloss.Style
			if alive {
				style = aliveStyle
			} else {
				style = deadStyle
			}

			if m.paused && y == m.posY && x == m.posX {
				line += style.Render("[]")
				continue
			}

			if alive {
				line += style.Render("  ")
			} else {
				line += style.Render("  ")
			}
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

type TickMsg struct{}

func (m Model) tickOnce() tea.Cmd {

	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}
