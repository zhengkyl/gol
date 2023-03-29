package life

// change to int or struct for customization
type Life bool

func NewBoard(width, height int) [][]Life {
	board := make([][]Life, height)
	for i := range board {
		board[i] = make([]Life, width)
	}

	return board
}

func NextBoard(board [][]Life) [][]Life {

	boardWidth := len(board[0])
	boardHeight := len(board)

	newBoard := NewBoard(boardWidth, boardHeight)

	for y := range board {
		for x := range board[y] {

			neighbors := 0

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					ny := (y + dy + boardHeight) % boardHeight
					nx := (x + dx + boardWidth) % boardWidth

					if board[ny][nx] {
						neighbors++
					}
				}
			}

			if !board[y][x] && neighbors == 3 {
				newBoard[y][x] = true
			}
			if board[y][x] && (neighbors == 2 || neighbors == 3) {
				newBoard[y][x] = true
			}
		}

	}
	return newBoard
}
