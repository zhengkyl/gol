package main

import "github.com/zhengkyl/gol/server"

// _ "net/http/pprof"

func main() {
	// go func() {
	// 	http.ListenAndServe("localhost:1234", nil)
	// }()
	server.RunServer()
}
