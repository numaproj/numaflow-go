TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/log-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/sink-log:${TAG}" --platform linux/amd64,linux/arm64 --target log . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/sink-log:${TAG}" --target log .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/sink-log:${TAG}"; fi

clean:
	-rm -rf ./dist
