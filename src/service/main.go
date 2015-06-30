package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"time"

	t "github.com/kihamo/go-wti"
	"github.com/vharitonsky/iniflags"
)

var (
	host                = flag.String("host", "0.0.0.0", "Service host")
	port                = flag.Int64("port", 9102, "Service port")
	wtiToken            = flag.String("wti_token", "", "WebTranslateIt API token")
	updateRetryDelay    = flag.Duration("update_retry_delay", 10*time.Second, "Update retry delay in seconds")
	updateRetryAttempts = flag.Int64("update_retry_attempts", 3, "Update retry attempts")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	iniflags.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	server, err := t.NewTranslatorServer(addr, *wtiToken, *updateRetryDelay, *updateRetryAttempts)
	if err != nil {
		log.Panic(err)
	}

	log.Print("Start server on ", addr)

	if err = server.Serve(); err != nil {
		log.Panic("Error calling serve on service: %s", err.Error())
	}
}
