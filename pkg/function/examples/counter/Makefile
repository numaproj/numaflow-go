.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/go-counter-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/go-counter-example:latest" --target counter .

clean:
	-rm -rf ./dist
