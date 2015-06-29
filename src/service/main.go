package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	t "github.com/kihamo/go-wti"
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
	addr := fmt.Sprintf("%s:%d", *host, *port)

	server, err := t.NewTranslatorServer(addr)
	if err != nil {
		log.Panic(err)
	}

	log.Print("Start server on ", addr)

	if err = server.Serve(); err != nil {
		log.Panic("Error calling serve on service: %s", err.Error())
	}
}
