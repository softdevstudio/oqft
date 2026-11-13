package oqft

import (
	"fmt"
	"math"
)

// optimisticQFT constructs Algorithm 1 of arXiv:2505.00701 with an
// explicit block width. Blocks are indexed in QFT traversal order.
func optimisticQFT(n, width int) (*Circuit, error) {
	if width < 1 {
		return nil, fmt.Errorf("QFT block width must be positive, got %d", width)
	}
	c, err := NewCircuit(n)
	if err != nil {
		return nil, err
	}

	type block struct {
		low, high int
		gates     []Gate
	}
	var blocks []block
	for high := n; high > 0; {
		// Anchor the partition at qubit 0, as in blockedQFT.
		low := (high - 1) / width * width
		local, err := QFT(high - low)
		if err != nil {
			return nil, err
		}
		// Retain the local bit reversal in the wire convention.
		gates := local.Gates[:local.Len()-(high-low)/2]
		for i := range gates {
			gates[i].Q0 += low
			if gates[i].Kind.Qubits() == 2 {
				gates[i].Q1 += low
			}
		}
		blocks = append(blocks, block{low, high, gates})
		high = low
	}

	phase := func(parity int) {
		for i := parity; i+1 < len(blocks); i += 2 {
			targets, controls := blocks[i], blocks[i+1]
			for target := targets.high - 1; target >= targets.low; target-- {
				for control := controls.high - 1; control >= controls.low; control-- {
					c.CP(control, target, math.Ldexp(math.Pi, control-target))
				}
			}
		}
	}

	for i := 0; i < len(blocks); i += 2 {
		c.Gates = append(c.Gates, blocks[i].gates...)
	}
	phase(0)
	for i := 1; i < len(blocks); i++ {
		gates := blocks[i].gates
		if i%2 == 1 {
			c.Gates = append(c.Gates, gates...)
			continue
		}
		// Undo the local QFT to estimate the input in the computational basis.
		for j := len(gates) - 1; j >= 0; j-- {
			gate := gates[j]
			gate.Angle = -gate.Angle
			c.Gates = append(c.Gates, gate)
		}
	}
	phase(1)
	for i := 2; i < len(blocks); i += 2 {
		c.Gates = append(c.Gates, blocks[i].gates...)
	}
	for low := 0; low < n/2; low++ {
		c.Swap(low, n-1-low)
	}
	return c, nil
}
