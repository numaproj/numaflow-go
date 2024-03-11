.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/sideinput-example main.go

.PHONY: image
image: build
	docker buildx build -t "quay.io/numaio/numaflow-go/sideinput-example:stable" --platform linux/amd64,linux/arm64 --target sideinput . --push

clean:
	-rm -rf ./dist
