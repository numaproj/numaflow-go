
.PHONY: test
test: proto generate
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


