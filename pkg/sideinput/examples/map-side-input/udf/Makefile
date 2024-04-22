DOCKERIO_ORG=quay.io/numaio
PLATFORM=linux/x86_64
TARGET=map-sideinput-udf

IMAGE_TAG=$(TAG)
ifeq ($(IMAGE_TAG),)
	IMAGE_TAG="dev"
endif

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/map-sideinput-udf main.go

.PHONY: image
#To build an image that supports multiple platforms(linux/amd64,linux/arm64) and push to quay.io, use the following command
image: build
	docker buildx build -t "$(DOCKERIO_ORG)/numaproj-contrib/e2e-map-sideinput-udf:$(IMAGE_TAG)" --platform $(PLATFORM) --target $(TARGET) . --load

.PHONY: lint
lint:
	go mod tidy
	golangci-lint run --fix --verbose --concurrency 4 --timeout 5m

clean:
	-rm -rf ./dist