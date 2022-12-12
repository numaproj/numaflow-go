.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/log-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/udsink/go-log:latest" --target log .

clean:
	-rm -rf ./dist
