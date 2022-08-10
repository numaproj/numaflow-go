//go:build tools

// This package contains code generation utilities
// This package imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import (
        _ "google.golang.org/protobuf/cmd/protoc-gen-go"
        _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
