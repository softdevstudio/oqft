package oqft

import (
	"fmt"
	"math/cmplx"
)

// Checker compares a circuit with the exact QFT on computational basis
// inputs. Use NewChecker to construct one.
type Checker struct {
	candidate Circuit
	exact     *Circuit
}

// NewChecker validates c and snapshots its gates for repeated checks.
func NewChecker(c *Circuit) (*Checker, error) {
	if c == nil {
		return nil, fmt.Errorf("checker needs a circuit")
	}
	if _, err := c.simulationSize(); err != nil {
		return nil, err
	}
	exact, err := QFT(c.N)
	if err != nil {
		return nil, err
	}
	return &Checker{candidate: Circuit{N: c.N, Gates: append([]Gate(nil), c.Gates...)}, exact: exact}, nil
}

// StateError returns the squared Euclidean distance between the two output
// vectors for basis. Averaging this over all basis indices gives the
// normalized squared Frobenius error.
func (c *Checker) StateError(basis int) (float64, error) {
	got, err := c.candidate.Simulate(basis)
	if err != nil {
		return 0, err
	}
	want, err := c.exact.Simulate(basis)
	if err != nil {
		return 0, err
	}
	var squaredError float64
	for i, amplitude := range got {
		delta := cmplx.Abs(amplitude - want[i])
		squaredError += delta * delta
	}
	return squaredError, nil
}
