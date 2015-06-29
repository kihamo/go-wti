package client

import (
	"log"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/kihamo/go-wti/gen-go/translator"
)

func NewTranslatorClient(addr string) (*translator.TranslatorClient, error) {
	var transport thrift.TTransport

	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		log.Fatal("Error starting client socket at %s: %s", addr, err)
	}

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport = transportFactory.GetTransport(transport)

	if err = transport.Open(); err != nil {
		return nil, err
	}

	return translator.NewTranslatorClientFactory(transport, protocolFactory), nil
}
