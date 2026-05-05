TAG ?= stable
PUSH ?= false
CONTAINER_RUNTIME ?= docker
IMAGE_REGISTRY = quay.io/numaio/numaflow-go/sink-no-response:${TAG}
ARCHITECTURES = amd64 arm64

.PHONY: build
build:
	for arch in $(ARCHITECTURES); do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -v -o ./dist/no-response-$${arch} main.go; \
	done

.PHONY: image-push
image-push: build
ifeq ($(CONTAINER_RUNTIME), docker)
	@echo "Using Docker to build and push the image..."
	docker buildx create --use --name no-response-builder 2>/dev/null || true
	docker buildx build --platform=linux/amd64,linux/arm64 --tag $(IMAGE_REGISTRY) --target no-response --push .
else ifeq ($(CONTAINER_RUNTIME), podman)
	@echo "Using Podman to build and push the image..."
	podman build --platform linux/amd64 --target no-response -t $(IMAGE_REGISTRY)-amd64 .
	podman build --platform linux/arm64 --target no-response -t $(IMAGE_REGISTRY)-arm64 .
	podman manifest create $(IMAGE_REGISTRY)
	podman manifest add $(IMAGE_REGISTRY) $(IMAGE_REGISTRY)-amd64
	podman manifest add $(IMAGE_REGISTRY) $(IMAGE_REGISTRY)-arm64
	podman manifest push --all $(IMAGE_REGISTRY)
else
	@echo "Unsupported container runtime: $(CONTAINER_RUNTIME)"
	@exit 1
endif

.PHONY: image
image: build
ifeq ($(CONTAINER_RUNTIME), docker)
	docker build -t ${IMAGE_REGISTRY} --target no-response .
	@if [ "$(PUSH)" = "true" ]; then docker push ${IMAGE_REGISTRY}; fi
else ifeq ($(CONTAINER_RUNTIME), podman)
	podman build -t ${IMAGE_REGISTRY} --target no-response .
	@if [ "$(PUSH)" = "true" ]; then podman push ${IMAGE_REGISTRY}; fi
else
	@echo "Unsupported container runtime: $(CONTAINER_RUNTIME)"
	@exit 1
endif

.PHONY: clean
clean:
	-rm -rf ./dist

