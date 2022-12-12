.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/counter-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/udf/reduce-go-counter:latest" --target counter .

clean:
	-rm -rf ./dist
