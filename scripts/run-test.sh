#!/usr/bin/env bash

function runTests {
  echo "==> Running Unit Tests..."
  go test -v $TEST "$TESTARGS" -timeout=600s -parallel=20
}

function main {
  runTests
}

main
