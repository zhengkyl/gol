package life

const DeadPlayer = 0

type Cell struct {
	Player       int
	PausedPlayer int
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
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					ny := (y + dy + boardHeight) % boardHeight
					nx := (x + dx + boardWidth) % boardWidth

					if board[ny][nx].Player != DeadPlayer {
						neighbors[board[ny][nx].Player]++
						numNeighbors++

						if neighbors[board[ny][nx].Player] > neighbors[mostColor] {
							mostColor = board[ny][nx].Player
						}
					}
				}
			}

			newBoard[y][x].PausedPlayer = board[y][x].PausedPlayer

			// One color must have a majority of 2 or 3 neighbors
			if neighbors[mostColor] < 2 {
				continue
			}

			if board[y][x].Player == DeadPlayer && numNeighbors == 3 {
				newBoard[y][x].Player = mostColor
			} else if board[y][x].Player != DeadPlayer && (numNeighbors == 2 || numNeighbors == 3) {
				newBoard[y][x].Player = mostColor
			}
		}

	}
	return newBoard
}
