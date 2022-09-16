.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/go-udsink-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/go-udsink-example:latest" --target flatmap .

clean:
	-rm -rf ./dist
