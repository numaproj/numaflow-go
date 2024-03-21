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

.PHONY: image-push-all
image-push-all:
	for dir in $(shell find ./pkg -name 'Makefile' -exec dirname {} \;); do \
		$(MAKE) -C $$dir image-push TAG=$(TAG); \
	done

