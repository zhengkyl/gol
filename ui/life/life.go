package life

// change to int or struct for customization

const deadColor = 0

type Cell struct {
	Color int
	Age   int
}

func (c *Cell) IsAlive() bool {
	return c.Color != deadColor
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
						neighbors[board[ny][nx].Color]++
						numNeighbors++

						if neighbors[board[ny][nx].Color] > mostNeighbors {
							mostColor = board[ny][nx].Color
						}
					}
				}
			}

			if !board[y][x].IsAlive() && numNeighbors == 3 {
				newBoard[y][x].Color = mostColor
			}
			if board[y][x].IsAlive() && (numNeighbors == 2 || numNeighbors == 3) {
				newBoard[y][x].Color = mostColor
			}
		}

	}
	return newBoard
}
