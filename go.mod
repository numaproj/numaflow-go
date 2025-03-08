module github.com/numaproj/numaflow-go

go 1.22

toolchain go1.23.1

require (
	github.com/stretchr/testify v1.9.0
	go.uber.org/atomic v1.11.0
	go.uber.org/mock v0.5.0
	golang.org/x/net v0.29.0
	golang.org/x/sync v0.8.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1
	google.golang.org/grpc v1.66.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.4.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v0.4.1
