.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/forward-message-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/map-forward-message:v0.5.0" --target forward-message .

clean:
	-rm -rf ./dist
