TAG ?= stable
PUSH ?= false
REPO = quay.io/numaio/numaflow-go/mapt-assign-event-time:${TAG}

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/assign-event-time-example main.go

.PHONY: image-push
image-push: build
	docker buildx build -t ${REPO} --platform linux/amd64,linux/arm64 --target assign-event-time . --push

.PHONY: image
image: build
	docker build -t ${REPO} --target assign-event-time .
	@if [ "$(PUSH)" = "true" ]; then docker push ${REPO}; fi

clean:
	-rm -rf ./dist
