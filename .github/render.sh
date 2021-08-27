#!/usr/bin/env bash

set -e
set -x
set -o pipefail

cd ./.github/render && \
  go build -o ../../render main.go && \
  cd ../../ && \
  ./render
