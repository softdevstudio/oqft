package oqft

import (
	"fmt"
	"math"
)

// GateKind is the type of a gate. The names match OpenQASM 3.
type GateKind int

const (
	// H is the Hadamard gate on one qubit.
	H GateKind = iota
	// CPhase is the controlled phase gate diag(1, 1, 1, e^{iθ}). It is
	// symmetric in its two qubits.
	CPhase
	// Swap exchanges two qubits.
	Swap
	// RZ is the Z rotation diag(e^{-iθ/2}, e^{iθ/2}) on one qubit.
	RZ
)

// Kinds lists every gate kind in a fixed order, for reports.
var Kinds = []GateKind{H, CPhase, Swap, RZ}

// String returns the OpenQASM 3 name of the gate kind.
func (k GateKind) String() string {
	switch k {
	case H:
		return "h"
	case CPhase:
		return "cp"
	case Swap:
		return "swap"
	case RZ:
		return "rz"
	}
	return fmt.Sprintf("GateKind(%d)", int(k))
}

// Qubits returns how many qubits a gate of this kind acts on, or 0 for an
// unknown kind.
func (k GateKind) Qubits() int {
	switch k {
	case H, RZ:
		return 1
	case CPhase, Swap:
		return 2
	}
	return 0
}

// HasAngle reports whether gates of this kind carry an angle.
func (k GateKind) HasAngle() bool { return k == CPhase || k == RZ }

// Gate is one gate in a circuit. Q1 is used by two-qubit gates only; for
// CPhase, Q0 is the control and Q1 the target. Angle is in radians and is
// used by CPhase and RZ only.
type Gate struct {
	Kind   GateKind
	Q0, Q1 int
	Angle  float64
}

// Circuit is a sequence of gates on N qubits, applied in order.
type Circuit struct {
	N     int
	Gates []Gate
}

// NewCircuit returns an empty circuit on n qubits.
func NewCircuit(n int) (*Circuit, error) {
	if n < 1 {
		return nil, fmt.Errorf("circuit needs at least 1 qubit, got %d", n)
	}
	return &Circuit{N: n}, nil
}

// Check reports whether g is a valid gate for a circuit on n qubits.
func (g Gate) Check(n int) error {
	q := g.Kind.Qubits()
	if q == 0 {
		return fmt.Errorf("unknown gate kind %d", int(g.Kind))
	}
	if g.Q0 < 0 || g.Q0 >= n {
		return fmt.Errorf("%s: qubit %d out of range 0..%d", g.Kind, g.Q0, n-1)
	}
	if q == 2 {
		if g.Q1 < 0 || g.Q1 >= n {
			return fmt.Errorf("%s: qubit %d out of range 0..%d", g.Kind, g.Q1, n-1)
		}
		if g.Q0 == g.Q1 {
			return fmt.Errorf("%s: both qubits are %d", g.Kind, g.Q0)
		}
	} else if g.Q1 != 0 {
		return fmt.Errorf("%s: one-qubit gate has a second qubit %d", g.Kind, g.Q1)
	}
	if g.Kind.HasAngle() {
		if math.IsNaN(g.Angle) || math.IsInf(g.Angle, 0) {
			return fmt.Errorf("%s: angle %v is not finite", g.Kind, g.Angle)
		}
	} else if g.Angle != 0 {
		return fmt.Errorf("%s: gate takes no angle, got %v", g.Kind, g.Angle)
	}
	return nil
}

// Add appends g after checking it. Use it for gates that come from
// outside the program.
func (c *Circuit) Add(g Gate) error {
	if err := g.Check(c.N); err != nil {
		return err
	}
	c.Gates = append(c.Gates, g)
	return nil
}

// must appends g and panics if it is invalid: in generator code a bad
// gate is a bug.
func (c *Circuit) must(g Gate) *Circuit {
	if err := c.Add(g); err != nil {
		panic("oqft: " + err.Error())
	}
	return c
}

// H appends a Hadamard on qubit q.
func (c *Circuit) H(q int) *Circuit { return c.must(Gate{Kind: H, Q0: q}) }

// CP appends a controlled phase of angle theta between control and target.
func (c *Circuit) CP(control, target int, theta float64) *Circuit {
	return c.must(Gate{Kind: CPhase, Q0: control, Q1: target, Angle: theta})
}

// Swap appends a swap of qubits a and b.
func (c *Circuit) Swap(a, b int) *Circuit { return c.must(Gate{Kind: Swap, Q0: a, Q1: b}) }

// RZ appends a Z rotation of angle theta on qubit q.
func (c *Circuit) RZ(q int, theta float64) *Circuit {
	return c.must(Gate{Kind: RZ, Q0: q, Angle: theta})
}

// Len returns the number of gates.
func (c *Circuit) Len() int { return len(c.Gates) }

// Counts returns the number of gates of each kind. Kinds with no gates are
// absent.
func (c *Circuit) Counts() map[GateKind]int {
	m := make(map[GateKind]int)
	for _, g := range c.Gates {
		m[g.Kind]++
	}
	return m
}

// TwoQubit returns the number of two-qubit gates.
func (c *Circuit) TwoQubit() int {
	n := 0
	for _, g := range c.Gates {
		if g.Kind.Qubits() == 2 {
			n++
		}
	}
	return n
}

// Depth returns the number of layers when every gate is placed as early
// as possible: a gate goes one layer after the latest earlier gate on any
// of its qubits. Every gate counts as one layer.
func (c *Circuit) Depth() int {
	top := make([]int, c.N)
	depth := 0
	for _, g := range c.Gates {
		l := top[g.Q0]
		if g.Kind.Qubits() == 2 && top[g.Q1] > l {
			l = top[g.Q1]
		}
		l++
		top[g.Q0] = l
		if g.Kind.Qubits() == 2 {
			top[g.Q1] = l
		}
		if l > depth {
			depth = l
		}
	}
	return depth
}
