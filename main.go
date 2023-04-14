package main

import (
	// _ "net/http/pprof"

	"github.com/zhengkyl/gol/server"
)

func main() {
	// go func() {
	// 	http.ListenAndServe("localhost:1234", nil)
	// }()
	server.RunServer()

	// uncomment to run ui without server
	// p := tea.NewProgram(ui.New(1, 1), tea.WithAltScreen())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("L + R, fix your code: %v", err)
	// 	os.Exit(1)
	// }
}
