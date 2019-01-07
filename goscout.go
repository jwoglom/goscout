package main

import (
	"flag"

	"./app"
	"github.com/ttacon/glog"
)

// DefaultPort is the default port for the server
const DefaultPort = 3000

var port = flag.Int("port", DefaultPort, "port to run server")
var testadd = flag.Bool("testadd", false, "test add")
var testget = flag.Bool("testget", false, "test get")

func main() {
	flag.Parse()

	s := app.NewServer(*port)

	if *testadd {
		glog.Infoln("adding fake treatments")
		s.Db.AddFakeTreatment()
		s.Db.AddFakeTreatment()
		s.Db.AddFakeTreatment()
	}

	if *testget {
	}

	s.Run()

}
