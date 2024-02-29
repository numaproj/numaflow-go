#!/bin/bash

commands=( "$@" )

find pkg -name "go.mod" | while read -r line;
do
    dir="$(dirname "${line}")"
    cd "$dir" || exit
    for cmd in "${commands[@]}"
        do
            eval "$cmd"
        done
    cd ~- || exit
done
