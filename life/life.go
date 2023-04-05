package life

func NewBoard(width, height int) [][]bool {
	board := make([][]bool, height)
	for i := range board {
		board[i] = make([]bool, width)
	}
	return board
}

func NextBoard(board [][]bool) [][]bool {
	boardHeight := len(board)
	boardWidth := len(board[0])

	newBoard := NewBoard(boardWidth, boardHeight)

	for y := range board {
		for x, alive := range board[y] {

			neighbors := 0

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dy == 0 && dx == 0 {
						continue
					}

					ny := (y + dy + boardHeight) % boardHeight
					nx := (x + dx + boardWidth) % boardWidth

					if board[ny][nx] {
						neighbors++
					}
				}
			}

			if alive && (neighbors == 2 || neighbors == 3) {
				newBoard[y][x] = true
			} else if !alive && neighbors == 3 {
				newBoard[y][x] = true
			}

		}
	}

	return newBoard
}
