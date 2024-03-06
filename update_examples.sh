#!/bin/bash

commit_sha=""

function show_help () {
    echo "Usage: $0 [-h|--help] (-b|--build | -c|--commit-sha <commit_sha>)"
    echo "  -h, --help         Display help message and exit"
    echo "  -b, --build        Build the docker images of all the examples"
    echo "  -c, --commit-sha   Update the numaflow-go version to the specified SHA"
}

function traverse_examples () {
  find pkg -name "go.mod" | while read -r line;
  do
      dir="$(dirname "${line}")"
      cd "$dir" || exit

      for command in "$@"
      do
        if ! $command; then
          echo "Error when running $command in $dir"
          exit 1
        fi
      done

      cd ~- || exit
  done
}

handle_options() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h | --help)
        show_help
        exit 0
        ;;
      -b | --build)
        traverse_examples "make clean"
        ;;
      -c | --commit-sha)
        if [ -z "$2" ]; then
          echo "Commit SHA not specified." >&2
          show_help
          exit 1
        fi

        commit_sha=$2
        traverse_examples "go get github.com/numaproj/numaflow-go@$commit_sha" "go mod tidy"
        shift
        ;;
      *)
        echo "Invalid option: $1" >&2
        show_help
        exit 1
        ;;
    esac
    shift
  done
}

echo "$#"
handle_options "$@"

if [ -n "$commit_sha" ]; then
 echo "Commit SHA specified: $commit_sha"
fi