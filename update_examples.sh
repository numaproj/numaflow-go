#!/bin/bash

function show_help () {
    echo "Usage: $0 [-h|--help] (-bp|--build-push | -bpe|--build-push-example | -u|--update <SDK-version>)"
    echo "  -h, --help                   Display help message and exit"
    echo "  -bp, --build-push            Build the Dockerfiles of all the examples and push them to the quay.io registry (with tag: stable)"
    echo "  -bpe, --build-push-example   Build the Dockerfile of the given example directory path, and push it to the quay.io registry (with tag: stable)"
    echo "  -u, --update                 Update all of the examples to depend on the numaflow-go version with the specified SHA or version"
}

function traverse_examples () {
  find pkg -name "go.mod" | while read -r line;
  do
      dir="$(dirname "${line}")"
      cd "$dir" || exit

      for command in "$@"
      do
        if ! $command; then
          echo "Error: failed $command in $dir" >&2
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
usingBuildPush=0
usingBuildPushExample=0
usingVersion=0
version=""
directoryPath=""

function handle_options () {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h | --help)
        usingHelp=1
        ;;
      -bp | --build-push)
        usingBuildPush=1
        ;;
      -bpe | --build-push-example)
        usingBuildPushExample=1
        if [ -z "$2" ]; then
          echo "Directory path not specified." >&2
          show_help
          exit 1
        fi

        directoryPath=$2
        shift
        ;;
      -u | --update)
        usingVersion=1
        if [ -z "$2" ]; then
          echo "Commit SHA or version not specified." >&2
          show_help
          exit 1
        fi

        version=$2
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

if (( usingBuildPush + usingBuildPushExample + usingVersion + usingHelp > 1 )); then
  echo "Only one of '-h', '-bp', '-bpe, or, '-u' is allowed at a time" >&2
  show_help
  exit 1
fi

if [ -n "$version" ]; then
 echo "Will update to: $version"
fi

if [ -n "$directoryPath" ]; then
 echo "Dockerfile path to use: $directoryPath"
fi

if (( usingBuildPush )); then
  traverse_examples "make image-push"
elif (( usingBuildPushExample )); then
   cd "./$directoryPath" || exit
   if ! make image-push; then
     echo "Error: failed to run make image in $directoryPath" >&2
     exit 1
   fi
elif (( usingVersion )); then
  traverse_examples "go get github.com/numaproj/numaflow-go@$version" "go mod tidy"
elif (( usingHelp )); then
  show_help
fi
