CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: all gen build clean precommit

all: clean build

gen:
	thrift -r --gen go:thrift_import="github.com/apache/thrift/lib/go/thrift" translator.thrift
	goimports -w $(CURRENT_DIR)gen-go

build: gen deps
	gb build

run-server: build
	$(CURRENT_DIR)bin/service -config=$(CURRENT_DIR)config/config.ini

run-client: build
    $(CURRENT_DIR)bin/client -config=$(CURRENT_DIR)config/config.ini

clean:
	rm -rf $(CURRENT_DIR)vendor/src/ $(CURRENT_DIR)bin $(CURRENT_DIR)pkg $(CURRENT_DIR)/gen-go
	go clean

deps:
	go get github.com/constabulary/gb/...
	gb vendor update -all

precommit:
	goimports -w $(CURRENT_DIR)
