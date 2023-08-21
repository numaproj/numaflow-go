.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/retry main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-retry:latest" --target retry .

clean:
	-rm -rf ./dist