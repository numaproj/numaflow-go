.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/even-odd-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-even-odd:latest" --target even-odd .

clean:
	-rm -rf ./dist
