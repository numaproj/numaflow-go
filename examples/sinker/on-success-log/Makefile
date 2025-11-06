TAG ?= stable
PUSH ?= false
IMAGE_REGISTRY = quay.io/numaio/numaflow-go/on-success-log:${TAG}
ARCHITECTURES = amd64 arm64

.PHONY: build
build:
	for arch in $(ARCHITECTURES); do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -v -o ./dist/on-success-log-$${arch} main.go; \
	done

.PHONY: image-push
image-push: build
	docker buildx build -t ${IMAGE_REGISTRY} --platform linux/amd64,linux/arm64 --target on_success_log . --push

.PHONY: image-load
image-load: build
	docker buildx build -t ${IMAGE_REGISTRY} --platform linux/arm64,linux/amd64 --target on_success_log . --load

.PHONY: image
image: build
	docker build -t ${IMAGE_REGISTRY} --target on_success_log .
	@if [ "$(PUSH)" = "true" ]; then docker push ${IMAGE_REGISTRY}; fi

clean:
	-rm -rf ./dist
