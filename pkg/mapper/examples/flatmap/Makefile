TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/flatmap-example main.go

.PHONY: image
image: build
	docker buildx build -t "quay.io/numaio/numaflow-go/map-flatmap:${TAG}" --platform linux/amd64,linux/arm64 --target flatmap . --push

.PHONY: test_image
test_image: build
	docker build -t "quay.io/numaio/numaflow-go/map-flatmap:${TAG}" --target flatmap .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/map-flatmap:${TAG}"; fi

clean:
	-rm -rf ./dist
