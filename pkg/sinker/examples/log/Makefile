.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/log-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/sink-log:v0.5.0" --target log .

clean:
	-rm -rf ./dist
