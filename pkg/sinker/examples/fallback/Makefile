TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/fallback-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/fb-sink-log:${TAG}" --platform linux/amd64,linux/arm64 --target fallback . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/fb-sink-log:${TAG}" --target fallback .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/fb-sink-log:${TAG}"; fi

clean:
	-rm -rf ./dist
