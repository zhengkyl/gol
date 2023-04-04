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

// 0 1 2 3     4 5 6 7     8
// 0 1 2 3 --- 0 1 2 3 --- 0
func NextBoard(board [][]Cell) [][]Cell {

	boardWidth := len(board[0])
	boardHeight := len(board)

	newBoard := NewBoard(boardWidth, boardHeight)

	for y := range board {
		for x := range board[y] {

			var neighbors uint64 = 0

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					ny := (y + dy + boardHeight) % boardHeight
					nx := (x + dx + boardWidth) % boardWidth

					if board[ny][nx].IsAlive() {
						neighbors++
					}
				}
			}

			if !board[y][x].IsAlive() && neighbors == 3 {
				newBoard[y][x].Color = 1
			}
			if board[y][x].IsAlive() && (neighbors == 2 || neighbors == 3) {
				newBoard[y][x].Color = 1
			}
		}

	}
	return newBoard
}
