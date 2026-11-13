# OQFT

OQFT is an open source toolkit for the **optimistic quantum Fourier
transform**: it builds the circuit, checks it on a classical simulator, and
works out what it saves when you count the cost of breaking RSA with
Shor's algorithm. Website: [oqft.com](https://oqft.com).

Download OQFT from the [latest release](https://github.com/softdevstudio/oqft/releases/latest).
See [CHANGELOG.md](CHANGELOG.md) for release notes.

## What it is, in plain words

The quantum Fourier transform (QFT) is the step in Shor's algorithm that
turns a hidden repeating pattern into a number you can read. Done exactly,
it's slow. In 2025 Kahanamoku-Meyer, Blue, Bergamaschi, Gidney and Chuang
showed it can be done much faster, in depth that grows with the logarithm
of the register size and with no spare qubits, if you accept that it's
wrong on a small slice of inputs that almost never come up
([arXiv:2505.00701](https://arxiv.org/abs/2505.00701)). They called these
*optimistic* circuits.

## Try it from source

Clone the repository and run these commands with Go 1.24 or newer:

```sh
go run ./cmd/oqft gen -n 4
go run ./cmd/oqft gen --type aqft -n 4 --cutoff 2
go run ./cmd/oqft stats -n 8 --cutoff 2
go run ./cmd/oqft gen --type oqft -n 128 --eps 0.01
go run ./cmd/oqft stats -n 128 --eps 0.01
```

`gen` writes the circuit as OpenQASM 3. `stats` compares the textbook
and approximate circuits, reporting gate counts by type, total gates,
two-qubit gates and depth. Each gate counts as one layer, and gates on
disjoint qubits can share a layer. The AQFT cutoff defaults to 2 and
keeps controlled rotations between qubits at most that far apart.

The Go library also simulates circuits and builds their full unitary
matrices. The [recorded results](test/results/) include sample stats and
the test run. Regenerate them with `scripts/record-results.sh`.

The optimistic generator uses epsilon as the bound
on normalized squared Frobenius error and defaults to 0.01. Small
registers may use the exact circuit when the bound calls for a width
that leaves at most two blocks. The RSA resource estimator and browser calculator are planned in
[CLAUDE.md](CLAUDE.md).

## Check a circuit

```sh
oqft check -n 4 --eps 0.01
oqft check -n 17 --eps 1 --samples 8 --seed 42
oqft check -n 12 --block-width 3 --basis 63 --threshold 0.5 --format csv
```

The checker compares output amplitudes with the textbook QFT for each
selected basis input. It checks all inputs by default through 10
qubits. For larger registers through 20 qubits, select one input with
`--basis` or draw distinct inputs with `--samples`. The seed makes a
sample reproducible. A sample mean describes only the selected inputs.

Epsilon controls circuit sizing. The independent per-state threshold
defaults to 0.1 and flags rows whose squared error exceeds it. Flagged
rows are reported as data, and a successfully written report exits zero.
For diagnostics, `--block-width` bypasses epsilon sizing and exposes
wraparound within the simulator limit. The
[recorded CSV](test/results/check-wraparound.csv) sweeps the estimated
block through all eight values with the neighboring block held fixed.

## Licence

Apache 2.0. See [LICENSE](LICENSE), [NOTICE](NOTICE) and
[THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md).

A [SoftDevStudio.com](https://softdevstudio.com) project.
