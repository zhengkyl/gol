package life

import (
	"math/rand"
	"testing"
)

func BenchmarkNextBoardDead(b *testing.B) {

	board := NewBoard(1000, 1000)
	for n := 0; n < b.N; n++ {
		board = NextBoard(board)
	}
}

func BenchmarkNextBoard10(b *testing.B) {

	board := make([][]Cell, 1000)
	for i := range board {
		board[i] = make([]Cell, 1000)
		for j := range board[i] {
			if rand.Intn(10) == 0 {
				board[i][j] = Cell{
					Player: 1,
				}
			}
		}
	}

	for n := 0; n < b.N; n++ {
		board = NextBoard(board)
	}
}
