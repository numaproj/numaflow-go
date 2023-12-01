SHELL:=/bin/bash

CURRENT_DIR=$(shell pwd)

.PHONY: all
all: proto generate test

clean:
	-rm -rf ${CURRENT_DIR}/dist

.PHONY: test
test:
	go test -race -v ./...

.PHONY: proto
proto: clean
	go mod vendor
	./hack/protogen.sh
	rm -rf ./vendor
	go mod tidy

.PHONY: generate
generate:
	go generate ./...


