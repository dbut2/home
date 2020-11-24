#!/bin/bash
set -e

if [ -z "$INPUT_IGNORE" ]; then
  # generate coverage profile
  GOFLAGS=-mod=vendor /go/bin/go-acc ./... -o=cover.out
else
  # generate coverage profile (with ignores)
  GOFLAGS=-mod=vendor /go/bin/go-acc ./... --ignore "${INPUT_IGNORE//[[:space:]]/}" -o=cover.out
fi

git config --global url."https://x-access-token:$INPUT_TOKEN@github.com".insteadOf "https://github.com"

GOFLAGS=-mod=vendor go tool cover -func=cover.out > /profile.txt
cd /coverage
cat /profile.txt | go run ./cmd/main/...
