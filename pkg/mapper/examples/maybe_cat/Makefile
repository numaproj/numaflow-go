TAG ?= stable
PUSH ?= false
IMAGE_REGISTRY = quay.io/numaio/numaflow-go/map-maybe-cat:${TAG}
ARCHITECTURES = amd64 arm64

.PHONY: build
build:
	for arch in $(ARCHITECTURES); do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -v -o ./dist/maybe-cat-example-$${arch} main.go; \
	done

.PHONY: image-push
image-push: build
	docker buildx build -t ${IMAGE_REGISTRY} --platform linux/amd64,linux/arm64 --target maybe-cat . --push

.PHONY: image
image: build
	docker build -t ${IMAGE_REGISTRY} --target maybe-cat .
	@if [ "$(PUSH)" = "true" ]; then docker push ${IMAGE_REGISTRY}; fi

clean:
	-rm -rf ./dist
