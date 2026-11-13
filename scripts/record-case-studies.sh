#!/usr/bin/env bash
# Records the runs behind the case studies and use examples on oqft.com to
# test/results/case-studies.txt. Every number in those posts comes from here.
set -euo pipefail
cd "$(dirname "$0")/.."
bin=$(mktemp -d)/oqft
go build -o "$bin" ./cmd/oqft
out=test/results/case-studies.txt
{
  echo "# Recorded $(date -u +%Y-%m-%d) at $(git rev-parse --short HEAD), $("$bin" version)"
  echo
  echo "## Depth and gate crossover, eps 0.01 (qubits qft_gates oqft_gates qft_depth oqft_depth)"
  for n in $(seq 32 80) 96 100 120 128; do
    "$bin" stats -n "$n" --eps 0.01 | awk -v n="$n" '/^Total gates/{qg=$3} /^Depth/{qd=$2} /^OQFT/{gsub(",","");print n, qg, $6, qd, $NF}'
  done
  echo
  echo "## Error targets at 2,048 qubits"
  for e in 0.1 0.01 0.001 0.0001; do
    go run ./scripts/blockwidths 2048 "$e"
    "$bin" stats -n 2048 --eps "$e" | grep '^OQFT'
  done
  echo
  echo "## Block width and measured error, 16 qubits, 300 inputs, seed 7, threshold 0.1"
  for w in 2 3 4 5 6 7 8; do
    echo "width $w: $("$bin" check -n 16 --block-width "$w" --samples 300 --seed 7 --threshold 0.1 | tail -1)"
  done
  echo
  echo "## Batch table (use example)"
  echo "qubits,eps,qft_depth,oqft_depth,qft_gates,oqft_gates"
  for n in 64 256 1024; do
    for eps in 0.1 0.01; do
      "$bin" stats -n "$n" --eps "$eps" | awk -v n="$n" -v e="$eps" '/^Total gates/{qg=$3} /^Depth/{qd=$2} /^OQFT/{gsub(",","");print n "," e "," qd "," $NF "," qg "," $6}'
    done
  done
} > "$out"
echo "recorded to $out"
