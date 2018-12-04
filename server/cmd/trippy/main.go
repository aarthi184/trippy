package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"trippy/server"

	_ "github.com/rakyll/gom/http"
)

func main() {

	// Live profiling webserver
	go func() {
		fmt.Println("See Profile info at http://localhost:6060/debug/pprof/ ")
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()

	var (
		err error
	)

	// Create New Server
	var s server.Service = new(server.Server)

	// Initialize
	if err = s.Initialize(); err != nil {
		fmt.Printf("Unable to initialize server. [Error:%s]\n", err)
		return
	}
	// Start the server
	if err = s.Start(); err != nil {
		fmt.Printf("Unable to start server. [Error:%s]\n", err)
		return
	}
}
