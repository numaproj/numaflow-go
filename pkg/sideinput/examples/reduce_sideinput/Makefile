TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/reduce-sideinput-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/reduce-sideinput:${TAG}" --platform linux/amd64,linux/arm64 --target sideinput . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/reduce-sideinput:${TAG}" --target sideinput .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/reduce-sideinput:${TAG}"; fi

clean:
	-rm -rf ./dist
