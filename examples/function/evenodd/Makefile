.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/go-even-odd-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/go-even-odd-example:latest" --target evenodd .

clean:
	-rm -rf ./dist
