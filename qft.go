package oqft

import "math"

// QFT returns the textbook quantum Fourier transform on n qubits.
// Its matrix element at output y and input x is exp(2πixy/2^n)/sqrt(2^n).
// The final swaps put the output in the same basis order as the input.
func QFT(n int) (*Circuit, error) {
	c, err := NewCircuit(n)
	if err != nil {
		return nil, err
	}
	for target := n - 1; target >= 0; target-- {
		c.H(target)
		for control := target - 1; control >= 0; control-- {
			angle := math.Ldexp(math.Pi, control-target)
			c.CP(control, target, angle)
		}
	}
	for low := 0; low < n/2; low++ {
		c.Swap(low, n-1-low)
	}
	return c, nil
}
