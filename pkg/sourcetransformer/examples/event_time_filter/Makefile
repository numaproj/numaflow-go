.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/event-time-filter-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/mapt-event-time-filter:stable" --target event-time-filter .

clean:
	-rm -rf ./dist


