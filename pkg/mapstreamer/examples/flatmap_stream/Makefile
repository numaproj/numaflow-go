TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/flatmap-stream-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/map-flatmap-stream:${TAG}" --platform linux/amd64,linux/arm64 --target flatmap_stream . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-flatmap-stream:${TAG}" --target flatmap_stream .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/map-flatmap-stream:${TAG}"; fi

clean:
	-rm -rf ./dist
