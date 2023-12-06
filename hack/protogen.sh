#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly REPO_ROOT="$(git rev-parse --show-toplevel)"

if [ "`command -v protoc`" = "" ]; then
  echo "Please install protobuf with - brew install protobuf"
  exit 1
fi

export PATH="$PATH:$(go env GOPATH)/bin"

if [ "`command -v protoc-gen-go`" = "" ]; then
  go install -mod=vendor ./vendor/google.golang.org/protobuf/cmd/protoc-gen-go
fi

if [ "`command -v protoc-gen-go-grpc`" = "" ]; then
  go install -mod=vendor ./vendor/google.golang.org/grpc/cmd/protoc-gen-go-grpc
fi

mkdir -p ${REPO_ROOT}/dist

ORG=${ORG:-numaproj}
PROJECT=${PROJECT:-numaflow}
BRANCH=${BRANCH:-main}

REMOTE_URL_PRE="https://raw.githubusercontent.com/${ORG}/${PROJECT}/${BRANCH}/pkg/apis/proto"
echo "Remote URL prefix: $REMOTE_URL_PRE"

curl -Ls -o ${REPO_ROOT}/dist/map.proto ${REMOTE_URL_PRE}/map/v1/map.proto
curl -Ls -o ${REPO_ROOT}/dist/reduce.proto ${REMOTE_URL_PRE}/reduce/v1/reduce.proto
curl -Ls -o ${REPO_ROOT}/dist/mapstream.proto ${REMOTE_URL_PRE}/mapstream/v1/mapstream.proto
curl -Ls -o ${REPO_ROOT}/dist/sideinput.proto ${REMOTE_URL_PRE}/sideinput/v1/sideinput.proto
curl -Ls -o ${REPO_ROOT}/dist/sink.proto ${REMOTE_URL_PRE}/sink/v1/sink.proto
curl -Ls -o ${REPO_ROOT}/dist/source.proto ${REMOTE_URL_PRE}/source/v1/source.proto
curl -Ls -o ${REPO_ROOT}/dist/transform.proto ${REMOTE_URL_PRE}/sourcetransform/v1/transform.proto

protoc --go_out=module=github.com/numaproj/numaflow-go:. --go-grpc_out=module=github.com/numaproj/numaflow-go:. -I ${REPO_ROOT}/dist $(find ${REPO_ROOT}/dist -name '*.proto')
