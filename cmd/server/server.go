package main

import (
	"flag"
	"fmt"

	"log"
	"time"

	"github.com/kihamo/godic"
)

var (
	host                = flag.String("host", "0.0.0.0", "Service host")
	port                = flag.Int64("port", 9102, "Service port")
	wtiToken            = flag.String("wti-token", "", "Webtranslateit api token")
	updateRetryDelay    = flag.Duration("update-retry-delay", 10*time.Second, "Update retry delay")
	updateRetryAttempts = flag.Int64("update-retry-attempts", 3, "Update retry delay")
	debug               = flag.Bool("debug", false, "Debug mode")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	server, err := godic.NewServer(addr, *debug)
	if err != nil {
		log.Fatal(err)
	}

	wti := godic.NewWebTranslateIt(*wtiToken, *updateRetryDelay, *updateRetryAttempts)
	server.SetWebTranslateIt(wti)

	log.Printf("turnpike server starting on %s\n", addr)
	log.Fatal(server.ListenAndServe())
}
