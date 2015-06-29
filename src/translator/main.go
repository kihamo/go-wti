package main

import (
	"log"
	"runtime"

	"github.com/kihamo/go-wti/service"
)

const (
	NetworkAddr = "0.0.0.0:9102"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	server, err := service.NewTranslatorServer(NetworkAddr)
	if err != nil {
		log.Panic(err)
	}

	if err = server.Serve(); err != nil {
		log.Panic("Error calling serve on hello server: %s", err.Error())
	}
}
