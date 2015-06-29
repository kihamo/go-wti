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
	addr := fmt.Sprintf("%s:%d", *host, *port)

	client, err := client.NewTranslatorClient(addr)
	if err != nil {
		log.Panic("Error start client: %s", err.Error())
	}
	defer client.Transport.Close()

	log.Print("Start client on ", addr)

	response, err := client.Ping()
	if err != nil {
		log.Fatal("Call method error: ", err)
	}
	fmt.Println(response)
	fmt.Println("ping()")
}
