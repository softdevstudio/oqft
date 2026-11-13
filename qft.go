package oqft

import (
	"fmt"
	"math"
)

// QFT returns the textbook quantum Fourier transform on n qubits.
// Its matrix element at output y and input x is exp(2πixy/2^n)/sqrt(2^n).
// The final swaps put the output in the same basis order as the input.
func QFT(n int) (*Circuit, error) {
	return qft(n, max(n, 1)-1)
}

// AQFT returns the approximate QFT on n qubits. It keeps controlled
// rotations whose qubits are at most cutoff positions apart. Zero keeps
// no controlled rotations. A cutoff of n-1 or more gives the textbook QFT.
func AQFT(n, cutoff int) (*Circuit, error) {
	if cutoff < 0 {
		return nil, fmt.Errorf("QFT cutoff must be nonnegative, got %d", cutoff)
	}
	return qft(n, cutoff)
}

func qft(n, cutoff int) (*Circuit, error) {
	c, err := NewCircuit(n)
	if err != nil {
		return nil, err
	}
	for target := n - 1; target >= 0; target-- {
		c.H(target)
		for control := target - 1; control >= max(0, target-cutoff); control-- {
			angle := math.Ldexp(math.Pi, control-target)
			c.CP(control, target, angle)
		}
	}
	for low := 0; low < n/2; low++ {
		c.Swap(low, n-1-low)
	}
	return c, nil
}
