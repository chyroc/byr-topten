#!/usr/bin/env bash

set -e
set -x
set -o pipefail

git config --global user.name 'bot'
git config --global user.email 'bot@github.com'
rm fetch || echo "fetch not exist"
git add .
git commit -am "commit-by-action: $(date)" || (echo "no commit" && exit 0)
git push