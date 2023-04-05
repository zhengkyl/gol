package main

import (
	"fmt"
	"gol/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.New(50, 25), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("L + R, fix your code: %v", err)
		os.Exit(1)
	}
}
