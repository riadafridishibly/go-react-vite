#!/bin/bash

set -e -x

pushd frontend
bun run build
popd

go build -ldflags '-s -w'