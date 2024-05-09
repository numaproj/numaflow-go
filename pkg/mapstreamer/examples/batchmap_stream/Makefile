TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/batchmap-stream-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/map-batchmap-stream:${TAG}" --platform linux/amd64,linux/arm64 --target batchmap_stream . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-batchmap-stream:${TAG}" --target batchmap_stream .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/map-batchmap-stream:${TAG}"; fi

clean:
	-rm -rf ./dist
