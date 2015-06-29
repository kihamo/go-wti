package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/kihamo/go-wti/service"
	"github.com/vharitonsky/iniflags"
)

var (
	host = flag.String("host", "0.0.0.0", "Service host")
	port = flag.Int64("port", 9102, "Service port")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	iniflags.Parse()

	server, err := service.NewTranslatorServer(fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Panic(err)
	}

	if err = server.Serve(); err != nil {
		log.Panic("Error calling serve on hello server: %s", err.Error())
	}
}
