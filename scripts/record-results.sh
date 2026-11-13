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
{
  for qubits in 4 8 16 32; do
    go run ./cmd/oqft stats -n "$qubits" --cutoff 2
    echo
  done
} > test/results/stats.txt
{
  for qubits in 32 64 128 256 512 1024 2048; do
    go run ./cmd/oqft stats -n "$qubits" --cutoff 2 --eps 0.01
    echo
  done
} > test/results/optimistic-stats.txt
go run ./cmd/oqft check -n 4 --eps 0.01 > test/results/check-exhaustive.txt
go run ./cmd/oqft check -n 17 --eps 1 --samples 8 --seed 42 > test/results/check-sampled.txt
{
  go run ./cmd/oqft check -n 12 --block-width 3 --basis 7 --threshold 0.5 --format csv
  for basis in 15 23 31 39 47 55 63; do
    go run ./cmd/oqft check -n 12 --block-width 3 --basis "$basis" --threshold 0.5 --format csv | sed '1d'
  done
} > test/results/check-wraparound.csv
echo "recorded to test/results"
