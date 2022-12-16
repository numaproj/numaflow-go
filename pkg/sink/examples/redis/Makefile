.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/redis-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/keranyang/numaflow/sink-go-redis:v1" --target redis .

clean:
	-rm -rf ./dist
