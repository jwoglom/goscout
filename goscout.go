package main

import (
	"flag"
	"math/rand"
	"time"

	"./app"
)

// DefaultPort is the default port for the server
const DefaultPort = 3000

var port = flag.Int("port", DefaultPort, "port to run server")
var testadd = flag.Bool("testadd", false, "test add")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	s := app.NewServer(*port)

	if *testadd {
		s.Db.AddFakeEntry()
	}

	s.Run()
}
