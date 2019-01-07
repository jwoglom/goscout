package main

import (
	"flag"

	"./app"
)

// DefaultPort is the default port for the server
const DefaultPort = 3000

var port = flag.Int("port", DefaultPort, "port to run server")

func main() {
	flag.Parse()

	s := app.NewServer(*port)
	s.Run()
}
