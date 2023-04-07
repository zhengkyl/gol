package middleware

import (
	"strings"
	"testing"

	"github.com/zhengkyl/gol/ui/life"
)

func BenchmarkNewlineIf(b *testing.B) {
	board := life.NewBoard(1000, 1000)

	for n := 0; n < b.N; n++ {
		sb := strings.Builder{}
		for y, row := range board {
			for _, cell := range row {
				style := deadStyle
				if cell.IsAlive() {
					style = aliveStyle
				}

				sb.WriteString(style.Render("  "))
			}

			if y != len(board)-1 {
				sb.WriteRune('\n')
			}
		}

		t := sb.String()
		_ = t
	}
}

func BenchmarkBuilderTechnique(b *testing.B) {
	board := life.NewBoard(1000, 1000)

	for n := 0; n < b.N; n++ {
		sb := strings.Builder{}

		for _, cell := range board[0] {
			style := deadStyle
			if cell.IsAlive() {
				style = aliveStyle
			}

			sb.WriteString(style.Render("  "))
		}
		for _, row := range board {
			for _, cell := range row {
				style := deadStyle
				if cell.IsAlive() {
					style = aliveStyle
				}

				sb.WriteString(style.Render("  "))
			}
			sb.WriteRune('\n')
		}

		t := sb.String()
		_ = t
	}
}

func BenchmarkSliceLast(b *testing.B) {
	board := life.NewBoard(1000, 1000)

	for n := 0; n < b.N; n++ {
		sb := strings.Builder{}

		for _, row := range board {
			for _, cell := range row {
				style := deadStyle
				if cell.IsAlive() {
					style = aliveStyle
				}

				sb.WriteString(style.Render("  "))
			}
			sb.WriteRune('\n')
		}

		t := sb.String()[:len(sb.String())-1]
		_ = t
	}
}

func BenchmarkJoin(b *testing.B) {

	board := life.NewBoard(1000, 1000)
	for n := 0; n < b.N; n++ {

		var lines []string
		for _, row := range board {
			sb := strings.Builder{}
			for _, cell := range row {
				style := deadStyle
				if cell.IsAlive() {
					style = aliveStyle
				}

				sb.WriteString(style.Render("  "))
			}

			lines = append(lines, sb.String())

		}
		t := strings.Join(lines, "\n")
		_ = t
	}
}
