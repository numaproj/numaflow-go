TAG ?= stable
PUSH ?= false

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/sum-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t "quay.io/numaio/numaflow-go/session-sum:${TAG}" --platform linux/amd64,linux/arm64 --target sum . --push

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/session-sum:${TAG}" --target sum .
	@if [ "$(PUSH)" = "true" ]; then docker push "quay.io/numaio/numaflow-go/session-sum:${TAG}"; fi

clean:
	-rm -rf ./dist
