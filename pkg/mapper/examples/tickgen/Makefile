TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/tickgen-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/map-tickgen:${TAG}" --platform linux/amd64,linux/arm64 --target tickgen . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-tickgen:${TAG}" --target tickgen .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/map-tickgen:${TAG}"; fi

clean:
	-rm -rf ./dist
