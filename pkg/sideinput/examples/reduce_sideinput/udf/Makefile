TAG ?= stable
PUSH ?= false
REPO = quay.io/numaio/numaflow-go/reduce-sideinput-udf:${TAG}

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/reduce-sideinput-udf main.go

.PHONY: image-push
image-push: build
	docker buildx build -t ${REPO} --platform linux/amd64,linux/arm64 --target sideinput-udf . --push

.PHONY: image
image: build
	docker build -t ${REPO} --target sideinput-udf .
	@if [ "$(PUSH)" = "true" ]; then docker push ${REPO}; fi

clean:
	-rm -rf ./dist
