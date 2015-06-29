package service

import (
	"log"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/kihamo/go-wti/gen-go/translator"
)

func NewTranslatorServer(addr string) (*thrift.TSimpleServer, error) {
	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Fatal("Error starting server socket at %s: %s", addr, err)
	}
	defer serverTransport.Close()

	handler := NewTranslatorHandler()
	processor := translator.NewTranslatorProcessor(handler)

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	return server, nil
}
