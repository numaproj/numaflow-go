.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/go-udsink-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/go-udsink-example:latest" --target udsink . && k3d image import quay.io/numaio/go-udsink-example:latest

clean:
	-rm -rf ./dist
