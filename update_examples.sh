#!/bin/bash

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
          echo "Error when running $command in $dir" >&2
          exit 1
        fi
      done

      cd ~- || exit
  done
}

if [ $# -eq 0 ]; then
  echo "Error: provide at least one argument" >&2
  show_help
  exit 1
fi

usingHelp=0
usingBuild=0
usingSHA=0
commit_sha=""

handle_options() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h | --help)
        usingHelp=1
        ;;
      -b | --build)
        usingBuild=1
        ;;
      -c | --commit-sha)
        usingSHA=1
        if [ -z "$2" ]; then
          echo "Commit SHA not specified." >&2
          show_help
          exit 1
        fi

        commit_sha=$2
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

handle_options "$@"

if (( usingBuild + usingSHA + usingHelp > 1 )); then
  echo "Only one of '-h', '-b', '-c' is allowed at a time" >&2
  show_help
  exit 1
fi

if [ -n "$commit_sha" ]; then
 echo "Commit SHA specified: $commit_sha"
fi

if (( usingBuild )); then
  traverse_examples "make image"
elif (( usingSHA )); then
  traverse_examples "go get github.com/numaproj/numaflow-go@$commit_sha" "go mod tidy"
elif (( usingHelp )); then
  show_help
fi
