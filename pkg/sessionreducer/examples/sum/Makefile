.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/sum-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/session-sum:v0.6.1" --target sum .

clean:
	-rm -rf ./dist
