package life

// change to int or struct for customization

const deadColor = 0

type Cell struct {
	Player       int
	PausedPlayer int
	Age          int
}

func (c *Cell) IsAlive() bool {
	return c.Player != deadColor
}

func NewBoard(width, height int) [][]Cell {
	board := make([][]Cell, height)
	for i := range board {
		board[i] = make([]Cell, width)
	}

	return board
}

func NextBoard(board [][]Cell) [][]Cell {

	boardWidth := len(board[0])
	boardHeight := len(board)

	newBoard := NewBoard(boardWidth, boardHeight)

	neighbors := map[int]int{}

	for y := range board {
		for x := range board[y] {

			numNeighbors := 0
			mostColor := 0
			mostNeighbors := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					ny := (y + dy + boardHeight) % boardHeight
					nx := (x + dx + boardWidth) % boardWidth

					if board[ny][nx].IsAlive() {
						neighbors[board[ny][nx].Player]++
						numNeighbors++

						if neighbors[board[ny][nx].Player] > mostNeighbors {
							mostColor = board[ny][nx].Player
						}
					}
				}
			}

			if !board[y][x].IsAlive() && numNeighbors == 3 {
				newBoard[y][x].Player = mostColor
			}
			if board[y][x].IsAlive() && (numNeighbors == 2 || numNeighbors == 3) {
				newBoard[y][x].Player = mostColor
				newBoard[y][x].Age = board[y][x].Age + 1
			}
		}

	}
	return newBoard
}
