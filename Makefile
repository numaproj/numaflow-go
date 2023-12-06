SHELL:=/bin/bash

CURRENT_DIR=$(shell pwd)

ORG?=numrproj
PROJECT?=numaflow
BRANCH?=main

.PHONY: all
all: generate test

clean:
	-rm -rf ${CURRENT_DIR}/dist

.PHONY: test
test:
	go test -race -v ./...

.PHONY: proto
proto: clean
	go mod vendor
	ORG=$(ORG) PROJECT=$(PROJECT) BRANCH=$(BRANCH) ./hack/protogen.sh
	rm -rf ./vendor
	go mod tidy

.PHONY: generate
generate:
	go generate ./...

