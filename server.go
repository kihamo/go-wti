package gowti

import (
	"log"

	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/kihamo/go-wti/gen-go/translator"
)

func NewTranslatorServer(addr string, wtiToken string, updateRetryDelay time.Duration, updateRetryAttempts int64) (*thrift.TSimpleServer, error) {
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Fatal("Error starting server socket at %s: %s", addr, err)
	}
	defer transport.Close()

	handler := NewTranslatorHandler(wtiToken, updateRetryDelay, updateRetryAttempts)
	processor := translator.NewTranslatorProcessor(handler)

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	return server, nil
}
