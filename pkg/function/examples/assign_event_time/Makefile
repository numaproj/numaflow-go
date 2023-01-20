.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/assign-event-time-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/mapt-assign-event-time:latest" --target assign-event-time .

clean:
	-rm -rf ./dist
