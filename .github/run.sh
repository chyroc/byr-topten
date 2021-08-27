#!/usr/bin/env bash

set -e
set -x
set -o pipefail

cd ./.github/cmd && \
  go build -o ../../dump main.go && \
  cd ../../ && \
  ./dump
