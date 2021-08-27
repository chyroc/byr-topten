#!/usr/bin/env bash

set -e
set -x
set -o pipefail

cd ./.github/fetch && \
  go build -o ../../fetch main.go && \
  cd ../../ && \
  ./fetch
