TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/map-sideinput-udf main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/map-sideinput-udf:${TAG}" --platform linux/amd64,linux/arm64 --target sideinput-udf . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-sideinput-udf:${TAG}" --target sideinput-udf .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/map-sideinput-udf:${TAG}"; fi

clean:
	-rm -rf ./dist
