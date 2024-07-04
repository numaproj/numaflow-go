TAG ?= stable
PUSH ?= false
IMAGE_REGISTRY = "quay.io/numaio/numaflow-go/redis-sink:${TAG}"

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/redis-sink main.go

.PHONY: image-push
image-push: build
	docker buildx build -t ${IMAGE_REGISTRY} --platform linux/amd64,linux/arm64 --target redis . --push

.PHONY: image
image: build
	docker build -t ${IMAGE_REGISTRY} --target redis .
	@if [ "$(PUSH)" = "true" ]; then docker push ${IMAGE_REGISTRY}; fi

clean:
	-rm -rf ./dist
