#!/bin/bash

go mod tidy
docker build . -t coverage-action
export INPUT_HARD_TARGET=40
export INPUT_SOFT_TARGET=95

cd ../../.. && \
docker run --rm \
  --workdir /github/workspace \
  -v $(pwd):/github/workspace \
  -e INPUT_HARD_TARGET \
  -e INPUT_SOFT_TARGET \
  -e ENABLE_SOFT_TARGET_WARNING="true" \
  -e INPUT_TOKEN="put github token here" \
  -e GITHUB_REPOSITORY="anzx/fabric-entitlements" \
  -e GITHUB_CONTEXT="{\"run_id\": \"163020168\", \"event\": {\"pull_request\": {\"number\": 431}}}" \
  coverage-action
