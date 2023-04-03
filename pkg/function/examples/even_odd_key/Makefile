.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/even-odd-key-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-even-odd-key:latest" --target even-odd-key .

clean:
	-rm -rf ./dist
