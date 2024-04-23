TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/redis-sink main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/redis-sink:${TAG}" --platform linux/amd64,linux/arm64 --target redis . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/redis-sink:${TAG}" --target redis .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/redis-sink:${TAG}"; fi

clean:
	-rm -rf ./dist
