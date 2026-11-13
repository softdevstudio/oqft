package oqft

import (
	"fmt"
	"strings"
)

// OpenQASM returns an OpenQASM 3 program for c, using the standard gate
// library and preserving gate order and qubit indices.
func (c *Circuit) OpenQASM() (string, error) {
	if c.N < 1 {
		return "", fmt.Errorf("circuit needs at least 1 qubit, got %d", c.N)
	}

	var program strings.Builder
	fmt.Fprintf(&program, "OPENQASM 3.0;\ninclude \"stdgates.inc\";\nqubit[%d] q;\n", c.N)
	for i, gate := range c.Gates {
		if err := gate.Check(c.N); err != nil {
			return "", fmt.Errorf("gate %d: %w", i, err)
		}
		program.WriteString(gate.Kind.String())
		if gate.Kind.HasAngle() {
			fmt.Fprintf(&program, "(%g)", gate.Angle)
		}
		fmt.Fprintf(&program, " q[%d]", gate.Q0)
		if gate.Kind.Qubits() == 2 {
			fmt.Fprintf(&program, ", q[%d]", gate.Q1)
		}
		program.WriteString(";\n")
	}
	return program.String(), nil
}
