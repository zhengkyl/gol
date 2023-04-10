package game

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/ui/life"
)

type PlayerState struct {
	Program *tea.Program
	PosX    int
	PosY    int
	Paused  bool
	Color   int
	Placed  int
	Cells   int
	Test    string
}

type GameState int

const (
	PAUSED GameState = iota
	PLAYING
)

type GameManager struct {
	games map[*Game]struct{}
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[*Game]struct{}),
	}
}

func (gm *GameManager) FindGame() *Game {
	for g := range gm.games {
		if g.players.Len() >= MaxPlayers {
			continue
		}
		return g
	}
	g := NewGame()
	g.Run()
	gm.games[g] = struct{}{}
	return g
}

func (gm *GameManager) EndGame(g *Game) {
	delete(gm.games, g)
}

type Game struct {
	players    PlayerMap
	board      [][]life.Cell
	boardMutex sync.Mutex
	paused     int
	ticker     *time.Ticker
	state      GameState
}

const MaxPlayers = 2
const MaxPlacedCells = 20
const drawRate = 10
const generationRate = 5
const drawsPerGeneration = drawRate / generationRate

const size = 100

func (g *Game) Players() int {
	return g.players.Len()
}

func (g *Game) PausedPlayers() int {
	return g.paused
}

func NewGame() *Game {
	w, h := size, size

	return &Game{
		players: PlayerMap{playerMap: make(map[int]*PlayerState)},
		board:   life.NewBoard(w, h),
		paused:  0,
		ticker:  time.NewTicker(time.Second / drawRate),
		state:   PAUSED,
	}
}

func (g *Game) Run() {
	go func() {

		var prevUpdate time.Time
		iteration := 0

		for now := range g.ticker.C {
			iteration++

			if iteration == drawsPerGeneration {
				iteration = 0
				g.UpdateBoard()
			}

			g.Update(now.Sub(prevUpdate))

			prevUpdate = now
		}
	}()
}

type JoinGameMsg struct {
	Game *Game
	Id   int
}

func (g *Game) Join(ps *PlayerState) (int, bool) {
	if g.players.Len() == MaxPlayers {
		return 0, false
	}

	posX := rand.Intn(len(g.board))
	posY := rand.Intn(len(g.board))

	ps.PosX = posX
	ps.PosY = posY
	ps.Paused = true
	ps.Color = g.players.Len() + 1

	id := g.players.Add(ps)

	return id, true
}

func (g *Game) Leave(id int) {
	g.boardMutex.Lock()
	defer g.boardMutex.Unlock()

	g.players.Delete(id)
	for y, row := range g.board {
		for x, cell := range row {
			if cell.Player == id {
				g.board[y][x].Player = 0
			}
		}
	}
}

func (g *Game) Unpause() {
	// update
}

func (g *Game) BoardSize() (int, int) {
	return len(g.board), len(g.board[0])
}

type ServerRedrawMsg struct{}

// type byPos []*ClientState

// func (s byPos) Len() int {
// 	return len(s)
// }
// func (s byPos) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }
// func (s byPos) Less(i, j int) bool {
// 	if s[i].PosX < s[j].PosX {
// 		return true
// 	}
// 	return s[i].PosY < s[j].PosY
// }

var deadStyle = lipgloss.NewStyle().Background(lipgloss.Color(ColorTable[0].cell))

func (g *Game) UpdateBoard() {

	if g.state != PLAYING {
		return
	}

	g.boardMutex.Lock()
	g.board = life.NextBoard(g.board)
	g.boardMutex.Unlock()
}

func (g *Game) Update(delta time.Duration) {
	g.paused = 0

	ids, pss := g.players.Entries()

	for _, cs := range pss {
		if cs.Paused {
			g.paused++
		}
	}

	if g.paused > len(ids)/2 {
		g.state = PAUSED
	} else {
		g.state = PLAYING
	}

	for _, p := range pss {
		p.Program.Send(ServerRedrawMsg{})
	}
}

func (g *Game) ViewBoard(top, left, width, height int) string {

	// Arbitrary limits to avoid unreasonable terminal sizes
	// This already shows the board 4 times
	if width > size*2 {
		width = size * 2
	}
	if height > size*2 {
		height = size * 2
	}

	sb := strings.Builder{}

	_, css := g.players.Entries()

	boardWidth, boardHeight := g.BoardSize()

	for y := top; y < top+height; y++ {
		boundY := (y + boardHeight) % boardHeight

		deadCount := 0
		for x := left; x < left+width; x++ {
			boundX := (x + boardWidth) % boardWidth
			style := lipgloss.NewStyle()
			pixel := "  "

			cursor := false
			for _, cs := range css {
				if boundY == cs.PosY && boundX == cs.PosX {
					cursor = true
					pixel = "[]"
					style = style.Background(lipgloss.Color(ColorTable[cs.Color].cursor))
				}
			}

			if !g.board[boundY][boundX].IsAlive() && !cursor {
				deadCount++
				continue
			}
			sb.WriteString(deadStyle.Render(strings.Repeat("  ", deadCount)))
			deadCount = 0

			if !cursor {
				style = style.Background(lipgloss.Color(ColorTable[g.board[boundY][boundX].Player].cell))
			}
			sb.WriteString(style.Render(pixel))
		}

		if deadCount > 0 {
			sb.WriteString(deadStyle.Render(strings.Repeat("  ", deadCount)))
		}

		sb.WriteString("\n")
	}

	return sb.String()[:sb.Len()-1]
}

func (g *Game) Place(id int) {

	cs := g.players.Get(id)

	// Maybe if player leaves but place() hasn't run yet?
	if cs == nil {
		return
	}

	if !g.board[cs.PosY][cs.PosX].IsAlive() {
		if cs.Placed >= MaxPlacedCells {
			return
		}
		g.board[cs.PosY][cs.PosX].Player = cs.Color
		cs.Placed++
	} else if g.board[cs.PosY][cs.PosX].Player == cs.Color {
		g.board[cs.PosY][cs.PosX].Player = 0
		cs.Placed--
	}
}

func (g *Game) GetPlayer(id int) *PlayerState {
	return g.players.Get(id)
}
