.PHONY: proto
proto:
	go mod vendor
	./hack/proto-gen.sh
	rm -rf ./vendor
	go mod tidy
