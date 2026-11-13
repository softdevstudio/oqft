# Changelog

## 0.2.0

- Added the ancilla-free optimistic QFT generator and OpenQASM 3 export
  with `oqft gen --type oqft`.
- Added optional OQFT statistics with `--eps` and recorded comparisons
  through 2,048 qubits.

## 0.1.0

- Added the circuit model, statevector simulator and full-unitary builder.
- Added textbook and approximate QFT generation and OpenQASM 3 export
  through `oqft gen`.
- Added `oqft stats` to compare gate counts and depth, and recorded example
  outputs alongside the tests.

## 0.0.0

- Repository setup: Go module, `oqft version`, CI, release workflow, and
  the shared instructions for Claude and ChatGPT/Codex. Not released.
