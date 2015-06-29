CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

all: clean build

build: deps
	thrift -r --gen go:thrift_import="github.com/apache/thrift/lib/go/thrift" translator.thrift
	gb build

run: build
	$(CURRENT_DIR)bin/translator

clean:
	rm -rf $(CURRENT_DIR)vendor/src/ $(CURRENT_DIR)/src/translator/gen-go bin pkg
	go clean

deps:
	go get github.com/constabulary/gb/...
	gb vendor update -all

precommit:
	goimports -w $(CURRENT_DIR)src/
