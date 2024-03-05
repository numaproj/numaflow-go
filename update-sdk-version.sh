#!/bin/bash

commit_sha=$1

find pkg -name "go.mod" | while read -r line;
do
    dir="$(dirname "${line}")"
    cd "$dir" || exit

    if ! go get github.com/numaproj/numaflow-go@"$commit_sha"; then
      echo "Failed to fetch latest SDK version in directory $dir"
      exit 1
    fi

    if ! go mod tidy; then
      echo "Failed to run go mod tidy in directory $dir"
      exit 1
    fi

    cd ~- || exit
done
