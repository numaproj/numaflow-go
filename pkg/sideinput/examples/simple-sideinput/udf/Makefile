.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/udf-sideinput-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/udf-sideinput:v3" --target udf-sideinput .

clean:
	-rm -rf ./dist
