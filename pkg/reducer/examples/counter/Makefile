TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/counter-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/reduce-counter:${TAG}" --platform linux/amd64,linux/arm64 --target counter . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/reduce-counter:${TAG}" --target counter .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/reduce-counter:${TAG}"; fi

clean:
	-rm -rf ./dist
