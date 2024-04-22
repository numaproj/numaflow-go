DOCKERIO_ORG=quay.io/numaio
PLATFORM=linux/amd64
TARGET=simple-source

IMAGE_TAG=$(TAG)
ifeq ($(IMAGE_TAG),)
    IMAGE_TAG=latest
endif

.PHONY: build image lint clean

build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/simple-source-example main.go

image: build
	docker buildx build --no-cache -t "$(DOCKERIO_ORG)/numaflow-source/simple_source_with_sideinput:$(IMAGE_TAG)" --platform $(PLATFORM) --target $(TARGET) . --load

lint:
	go mod tidy
	golangci-lint run --fix --verbose --concurrency 4 --timeout 5m

clean:
	-rm -rf ./dist
