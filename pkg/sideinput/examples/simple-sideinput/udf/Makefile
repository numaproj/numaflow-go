TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/udf-sideinput-example main.go

.PHONY: image
image: build
	docker buildx build -t "quay.io/numaio/numaflow-go/udf-sideinput-example:${TAG}" --platform linux/amd64,linux/arm64 --target udf-sideinput . --push

# If testing locally, there isn't a need for multi-arch images
.PHONY: test_image
test_image: build
	docker build -t "quay.io/numaio/numaflow-go/udf-sideinput-example:${TAG}" --target udf-sideinput .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/udf-sideinput-example:${TAG}"; fi

clean:
	-rm -rf ./dist
