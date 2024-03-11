.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/udf-sideinput-example main.go

.PHONY: image
#To build an image that supports multiple platforms(linux/amd64,linux/arm64) and push to quay.io, use the following command
image: build
	docker buildx build -t "quay.io/numaio/numaflow-go/udf-sideinput-example:stable" --platform linux/amd64,linux/arm64 --target udf-sideinput .



clean:
	-rm -rf ./dist
