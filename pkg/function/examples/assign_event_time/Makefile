.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/assign-event-time-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/assign-event-time:latest" --target assign-event-time-example .

clean:
	-rm -rf ./dist
