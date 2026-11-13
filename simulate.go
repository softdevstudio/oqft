package oqft

import (
	"fmt"
	"math"
	"math/cmplx"
)

// MaxSimulationQubits bounds the memory needed for a statevector.
const MaxSimulationQubits = 20

// Simulate evolves a computational basis state through c and returns the
// resulting amplitudes. Qubit 0 is the least significant bit of basis.
func (c *Circuit) Simulate(basis int) ([]complex128, error) {
	size, err := c.simulationSize()
	if err != nil {
		return nil, err
	}
	if basis < 0 || basis >= size {
		return nil, fmt.Errorf("basis index %d out of range 0..%d", basis, size-1)
	}
	state := make([]complex128, size)
	state[basis] = 1
	c.evolve(state)
	return state, nil
}

// Evolve applies c in place to a vector of 2^N amplitudes in basis-index
// order. It preserves the input norm, so the vector need not be normalized.
// If validation fails, state is unchanged.
func (c *Circuit) Evolve(state []complex128) error {
	size, err := c.simulationSize()
	if err != nil {
		return err
	}
	if len(state) != size {
		return fmt.Errorf("state has %d amplitudes, want %d", len(state), size)
	}
	c.evolve(state)
	return nil
}

func (c *Circuit) simulationSize() (int, error) {
	if c.N < 1 || c.N > MaxSimulationQubits {
		return 0, fmt.Errorf("simulator supports 1..%d qubits, got %d", MaxSimulationQubits, c.N)
	}
	for i, gate := range c.Gates {
		if err := gate.Check(c.N); err != nil {
			return 0, fmt.Errorf("gate %d: %w", i, err)
		}
	}
	return 1 << c.N, nil
}

func (c *Circuit) evolve(state []complex128) {
	for _, gate := range c.Gates {
		first := 1 << gate.Q0
		switch gate.Kind {
		case H:
			for block := 0; block < len(state); block += 2 * first {
				for i := block; i < block+first; i++ {
					j := i + first
					a, b := state[i], state[j]
					state[i] = (a + b) * (math.Sqrt2 / 2)
					state[j] = (a - b) * (math.Sqrt2 / 2)
				}
			}
		case CPhase:
			both := first | (1 << gate.Q1)
			phase := cmplx.Rect(1, gate.Angle)
			for i := range state {
				if i&both == both {
					state[i] *= phase
				}
			}
		case Swap:
			both := first | (1 << gate.Q1)
			for i := range state {
				if i&both == first {
					j := i ^ both
					state[i], state[j] = state[j], state[i]
				}
			}
		case RZ:
			zero := cmplx.Rect(1, -gate.Angle/2)
			one := cmplx.Conj(zero)
			for block := 0; block < len(state); block += 2 * first {
				for i := block; i < block+first; i++ {
					state[i] *= zero
					state[i+first] *= one
				}
			}
		}
	}
}
