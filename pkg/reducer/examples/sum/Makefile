.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/sum-example main.go

.PHONY: image
image: build
	docker build -t "quay.io/numaio/numaflow-go/reduce-sum:stable" --target sum .

clean:
	-rm -rf ./dist
