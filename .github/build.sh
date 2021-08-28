#!/usr/bin/env bash

set -e
set -x
set -o pipefail

go install github.com/chyroc/action.sh/commiter@v0.4.0 && mv `which commiter` /tmp/commiter
( cd ./.github/cmd-fetch  && go build -o /tmp/fetch  main.go )
( cd ./.github/cmd-render && go build -o /tmp/render main.go )
