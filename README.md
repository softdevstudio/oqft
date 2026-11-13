# OQFT

OQFT is an open source toolkit for the **optimistic quantum Fourier
transform**: it builds the circuit, checks it on a classical simulator, and
works out what it saves when you count the cost of breaking RSA with
Shor's algorithm. Website: [oqft.com](https://oqft.com).

> Status: setting up. The first working release is 0.1.0. See
> [CHANGELOG.md](CHANGELOG.md).

## What it is, in plain words

The quantum Fourier transform (QFT) is the step in Shor's algorithm that
turns a hidden repeating pattern into a number you can read. Done exactly,
it's slow. In 2025 Kahanamoku-Meyer, Blue, Bergamaschi, Gidney and Chuang
showed it can be done much faster, in depth that grows with the logarithm
of the register size and with no spare qubits, if you accept that it's
wrong on a small slice of inputs that almost never come up
([arXiv:2505.00701](https://arxiv.org/abs/2505.00701)). They called these
*optimistic* circuits.

OQFT turns that paper into a tool you can run:

- `oqft gen`: the textbook, approximate and optimistic QFT as OpenQASM 3,
  for any register size and error target
- `oqft stats`: gate counts and depth, side by side
- `oqft check`: runs the circuit on a simulator and shows where it fails
- `oqft estimate`: qubits and depth for factoring an RSA key, with and
  without the optimistic QFT

The same estimator runs in the browser as the Q-Day calculator on
oqft.com.

## Install

Download the binary for your system from the
[latest release](https://github.com/softdevstudio/oqft/releases/latest).
Or, with Go installed:

```
go install github.com/softdevstudio/oqft/cmd/oqft@latest
```

## Licence

Apache 2.0. See [LICENSE](LICENSE), [NOTICE](NOTICE) and
[THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md).

A [SoftDevStudio.com](https://softdevstudio.com) project.
