#!/usr/bin/env bash
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}")" && pwd )

pushd "${DIR}"/.. || exit
go build -o example/protoc-gen-gsk .
popd || exit
