TAG ?= stable
PUSH ?= false
REPO = quay.io/numaio/numaflow-go/mapt-event-time-filter:${TAG}

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/event-time-filter-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t ${REPO} --platform linux/amd64,linux/arm64 --target event-time-filter . --push

.PHONY: image
image: build
	docker build -t ${REPO} --target event-time-filter .
	@if [ "$(PUSH)" = "true" ]; then docker push ${REPO}; fi

clean:
	-rm -rf ./dist


