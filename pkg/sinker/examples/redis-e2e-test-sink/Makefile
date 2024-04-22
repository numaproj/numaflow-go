.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/redis-e2e-test-sink main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-sink/redis-e2e-test-sink:latest" --target redis .

clean:
	-rm -rf ./dist