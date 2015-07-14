package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kihamo/godic"
)

var (
	host     = flag.String("host", "0.0.0.0", "Service host")
	port     = flag.Int64("port", 9102, "Service port")
	wtiToken = flag.String("wti-token", "", "Webtranslateit api token")
	debug    = flag.Bool("debug", false, "Debug mode")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	wti := godic.NewWebTranslateIt(*wtiToken)
	server, err := godic.NewServer(addr, wti, *debug)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("turnpike server starting on %s\n", addr)
	log.Fatal(server.ListenAndServe())
}
