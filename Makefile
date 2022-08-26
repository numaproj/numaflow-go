.PHONY: proto

proto:
	go mod vendor
	./hack/protogen.sh
	rm -rf ./vendor
	go mod tidy

generate:
	go generate ./...