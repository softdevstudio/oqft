package oqft

import (
	"fmt"
)

// MaxUnitaryQubits bounds the memory needed for a full unitary matrix.
const MaxUnitaryQubits = 10

// Unitary returns the circuit's matrix indexed as matrix[output][input].
// Each column is the result of evolving the corresponding basis state.
func (c *Circuit) Unitary() ([][]complex128, error) {
	if c.N > MaxUnitaryQubits {
		return nil, fmt.Errorf("unitary builder supports up to %d qubits, got %d", MaxUnitaryQubits, c.N)
	}
	size, err := c.simulationSize()
	if err != nil {
		return nil, err
	}

	matrix := make([][]complex128, size)
	for output := range matrix {
		matrix[output] = make([]complex128, size)
	}
	column := make([]complex128, size)
	for input := 0; input < size; input++ {
		clear(column)
		column[input] = 1
		c.evolve(column)
		for output, amplitude := range column {
			matrix[output][input] = amplitude
		}
	}
	return matrix, nil
}
