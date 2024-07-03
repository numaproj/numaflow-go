TAG ?= stable
PUSH ?= false
REPO = quay.io/numaio/numaflow-go/fb-sink-log:${TAG}

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/fallback-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t ${REPO} --platform linux/amd64,linux/arm64 --target fallback . --push

.PHONY: image
image: build
	docker build -t ${REPO} --target fallback .
	@if [ "$(PUSH)" = "true" ]; then docker push ${REPO}; fi

clean:
	-rm -rf ./dist
