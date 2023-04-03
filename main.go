package main

import "github.com/zhengkyl/gol/server"

func main() {
	server.RunServer()

	// uncomment to run ui without server
	// p := tea.NewProgram(ui.New(1, 1), tea.WithAltScreen())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("L + R, fix your code: %v", err)
	// 	os.Exit(1)
	// }
}
