package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kihamo/go-wti/client"
	"github.com/vharitonsky/iniflags"
)

var (
	host = flag.String("host", "0.0.0.0", "Service host")
	port = flag.Int64("port", 9102, "Service port")
)

func main() {
	iniflags.Parse()

	client, err := client.NewTranslatorClient(fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Panic(err)
	}

	err = client.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ping()")
}
