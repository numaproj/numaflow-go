.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/flatmap-stream-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-flatmap-stream:v0.5.0" --target flatmap_stream .

clean:
	-rm -rf ./dist
