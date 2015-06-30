package main

import (
	"flag"
	"fmt"
	"log"

	t "github.com/kihamo/go-wti"
	"github.com/vharitonsky/iniflags"
)

var (
	host = flag.String("host", "0.0.0.0", "Service host")
	port = flag.Int64("port", 9102, "Service port")
)

func main() {
	iniflags.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	client, err := t.NewTranslatorClient(addr)
	if err != nil {
		log.Panic(err)
	}
	defer client.Transport.Close()

	log.Print("Start client on ", addr)

	response, err := client.GetDictionary("id_id1")
	if err != nil {
		log.Fatal("Call method error: ", err)
	}
	fmt.Println(response["Payment Methods"])
}
