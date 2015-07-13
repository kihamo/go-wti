package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kihamo/godic/sdk"
)

var (
	host  = flag.String("host", "0.0.0.0", "Service host")
	port  = flag.Int64("port", 9102, "Service port")
	debug = flag.Bool("debug", false, "Debug mode")
)

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	client, err := sdk.NewClient(addr, *debug)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan bool)
	client.UpdateSubscribe(func() {
		log.Println("Local update")
		//quit <- true
	})

	/*
		update, err := client.DictionaryUpdate()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send dictionary update result %v\n", update)
	*/

	/*
		dic, err := client.GetDictionary("vi_vn")
		if err != nil {
			log.Fatal(err)
		}
	*/

	<-quit
}
