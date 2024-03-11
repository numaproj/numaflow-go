TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/assign-event-time-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/mapt-assign-event-time:${TAG}" --platform linux/amd64,linux/arm64 --target assign-event-time . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/mapt-assign-event-time:${TAG}" --target assign-event-time .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/mapt-assign-event-time:${TAG}"; fi

clean:
	-rm -rf ./dist
