package game

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/ui/life"
)

type ClientState struct {
	PosX int
	PosY int
	// Pos    Coord
	Paused bool
	Color  int
}
type pixelLookup struct {
	start int
	end   int
}
type Game struct {
	clients      map[*tea.Program]*ClientState
	board        [][]life.Cell
	bufferLookup [][]pixelLookup
	buffer       string
	players      int
	ticker       *time.Ticker
}

type playerColor struct {
	cursor string
	cell   string
}

var ColorTable = [11]playerColor{
	{
		"#080808",
		"0",
	},
	{
		"#ff0000",
		"#ff5f5f",
	},
	{
		"#d75f00",
		"#ff8700",
	},
	{
		"#ffd700",
		"#ffff5f",
	},
	{
		"#87af00",
		"#afff00",
	},
	{
		"#005f00",
		"#00d700",
	},
	{
		"#00afff",
		"#00ffff",
	},
	{
		"#005f87",
		"#0087ff",
	},
	{
		"#d700ff",
		"#d787ff",
	},
	{
		"#ff00af",
		"#ff5faf",
	},
	{
		"#afafd7",
		"#eeeeee",
	},
}

const tickRate = 60
const drawRate = 10
const ticksPerDraw = tickRate / drawRate

func NewGame() *Game {

	w, h := 20, 20
	lookup := make([][]pixelLookup, h)
	for y := range lookup {
		lookup[y] = make([]pixelLookup, w)
	}

	return &Game{
		clients:      make(map[*tea.Program]*ClientState),
		board:        life.NewBoard(w, h),
		bufferLookup: lookup,
		players:      0,
		ticker:       time.NewTicker(time.Second / tickRate),
	}
}

func (g *Game) Run() {
	go func() {

		var prevUpdate time.Time
		iteration := 0

		for now := range g.ticker.C {
			iteration++

			if iteration == ticksPerDraw {
				iteration = 0
				g.UpdateBoard()
			}

			g.Update(now.Sub(prevUpdate))

			prevUpdate = now
		}
	}()
}

func (g *Game) Join(p *tea.Program, cs *ClientState) {
	posX := rand.Intn(len(g.board))
	posY := rand.Intn(len(g.board))
	g.players++

	cs.PosX = posX
	cs.PosY = posY
	cs.Paused = true
	cs.Color = g.players
	g.clients[p] = cs

}

func (g *Game) Leave(p *tea.Program) {
	delete(g.clients, p)
	g.players--
}

func (g *Game) Unpause() {
	// update
}

func (g *Game) Iterate() {
	// update
}

func (g *Game) BoardSize() (int, int) {
	return len(g.board), len(g.board[0])
}

type ServerRedrawMsg struct{}

type byPos []*ClientState

func (s byPos) Len() int {
	return len(s)
}
func (s byPos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPos) Less(i, j int) bool {
	if s[i].PosX < s[j].PosX {
		return true
	}
	return s[i].PosY < s[j].PosY
}

// var aliveStyle = lipgloss.NewStyle().Background(lipgloss.Color("227"))
var cornerStyle = lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[0].cursor))
var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[0].cell))

// var aliveWidth = len(aliveStyle.Render(pixel))

func (g *Game) UpdateBoard() {

	notPaused := true

	var clients []*ClientState
	for _, cs := range g.clients {
		clients = append(clients, cs)
		if cs.Paused {
			notPaused = false
		}
	}
	sort.Sort(byPos(clients))

	if notPaused {
		g.board = life.NextBoard(g.board)
	}

	sb := strings.Builder{}

	for y, row := range g.board {
		// line := ""
		for x, cell := range row {
			style := deadStyle
			if cell.IsAlive() {
				// style = lipgloss.NewStyle().Background(lipgloss.Color(strconv.Itoa(cell.Color)))
				style = lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[cell.Color].cell))
			} else if y == 0 && x == 0 {
				style = cornerStyle
			}
			var pixel = "  "

			for _, cs := range clients {
				if y == cs.PosY && x == cs.PosX {
					pixel = "[]"
					style = style.Copy().Foreground(lipgloss.Color(ColorTable[cs.Color].cursor))
				}
			}

			g.bufferLookup[y][x].start = sb.Len()
			sb.WriteString(style.Render(pixel))
			g.bufferLookup[y][x].end = sb.Len()
		}
		sb.WriteString("\n")
	}
	g.buffer = sb.String()[:sb.Len()-1]

	for p := range g.clients {
		p.Send(ServerRedrawMsg{})
	}
}

func (g *Game) ViewBoard(top, left, width, height int) string {
	sb := strings.Builder{}

	boardWidth, boardHeight := g.BoardSize()

	for y := top; y < top+height; y++ {

		boundY := (y + boardHeight) % boardHeight

		boundXStart := (left + boardWidth) % boardWidth
		boundXEndIncl := (left + width - 1 + boardWidth) % boardWidth

		// if width % boardWidth == 0 -> special case
		repeats := (width - 1) / boardWidth

		wrap := boundXStart > boundXEndIncl || repeats > 0

		// remove 1 when repeat is discontinuous
		if boundXStart <= boundXEndIncl && repeats > 0 {
			repeats--
		}

		if wrap {

			start := g.bufferLookup[boundY][boundXStart].start
			end := g.bufferLookup[boundY][boardWidth-1].end
			sb.WriteString(g.buffer[start:end])

			if repeats > 0 {
				start = g.bufferLookup[boundY][0].start
				end = g.bufferLookup[boundY][boardWidth-1].end
				sb.WriteString(strings.Repeat(g.buffer[start:end], repeats))
			}

			start = g.bufferLookup[boundY][0].start
			end = g.bufferLookup[boundY][boundXEndIncl].end
			sb.WriteString(g.buffer[start:end])

		} else {
			start := g.bufferLookup[boundY][boundXStart].start
			end := g.bufferLookup[boundY][boundXEndIncl].end
			sb.WriteString(g.buffer[start:end])
		}

		sb.WriteString("\n")
	}

	return sb.String()[:sb.Len()-1]
}

func (g *Game) Update(delta time.Duration) {

}

func (g *Game) Place(cs *ClientState) {
	if !g.board[cs.PosY][cs.PosX].IsAlive() {
		g.board[cs.PosY][cs.PosX].Color = cs.Color
	} else if g.board[cs.PosY][cs.PosX].Color == cs.Color {
		g.board[cs.PosY][cs.PosX].Color = 0
	}
}
