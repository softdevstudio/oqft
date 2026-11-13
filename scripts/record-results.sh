#!/usr/bin/env bash
# Records the test run and the environment to test/results. Every number on
# oqft.com comes from these files: re-run after any change that could move
# one, commit the results, then update the site from them.
set -euo pipefail
cd "$(dirname "$0")/.."
mkdir -p test/results
{
  echo "date: $(date -u +%Y-%m-%d)"
  echo "commit: $(git rev-parse --short HEAD 2>/dev/null || echo none)"
  echo "oqft: $(tr -d '[:space:]' < VERSION)"
  go version
  echo "os: $(uname -sm)"
} > test/results/environment.txt
go test -count=1 -v ./... 2>&1 | tee test/results/tests.txt
echo "recorded to test/results"
