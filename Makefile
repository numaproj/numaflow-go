.PHONY: all
all: proto generate test

.PHONY: test
test:
	go test -race -v ./...

.PHONY: proto
proto:
	go mod vendor
	./hack/protogen.sh
	rm -rf ./vendor
	go mod tidy

.PHONY: generate
generate:
	go generate ./...
