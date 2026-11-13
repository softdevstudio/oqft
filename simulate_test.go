package oqft

import (
	"math"
	"math/cmplx"
	"slices"
	"testing"
)

func TestSimulate(t *testing.T) {
	const amplitude = math.Sqrt2 / 2
	tests := []struct {
		name  string
		gates []Gate
		basis int
		want  []complex128
	}{
		{"identity", nil, 2, []complex128{0, 0, 1, 0}},
		{"Hadamard on low bit", []Gate{{Kind: H}}, 1, []complex128{amplitude, -amplitude, 0, 0}},
		{"Hadamard on high bit", []Gate{{Kind: H, Q0: 1}}, 0, []complex128{amplitude, 0, amplitude, 0}},
		{"controlled phase active", []Gate{{Kind: CPhase, Q1: 1, Angle: math.Pi / 2}}, 3, []complex128{0, 0, 0, 1i}},
		{"controlled phase inactive", []Gate{{Kind: CPhase, Q1: 1, Angle: math.Pi / 2}}, 2, []complex128{0, 0, 1, 0}},
		{"swap low to high", []Gate{{Kind: Swap, Q1: 1}}, 1, []complex128{0, 0, 1, 0}},
		{"swap high to low", []Gate{{Kind: Swap, Q1: 1}}, 2, []complex128{0, 1, 0, 0}},
		{"swap equal bits", []Gate{{Kind: Swap, Q1: 1}}, 3, []complex128{0, 0, 0, 1}},
		{"RZ on zero bit", []Gate{{Kind: RZ, Q0: 1, Angle: math.Pi}}, 1, []complex128{0, -1i, 0, 0}},
		{"RZ on one bit", []Gate{{Kind: RZ, Q0: 1, Angle: math.Pi}}, 2, []complex128{0, 0, 1i, 0}},
		{"controlled bit flip", []Gate{{Kind: H}, {Kind: CPhase, Q1: 1, Angle: math.Pi}, {Kind: H}}, 2, []complex128{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Circuit{N: 2, Gates: tt.gates}
			got, err := c.Simulate(tt.basis)
			if err != nil {
				t.Fatal(err)
			}
			checkAmplitudes(t, got, tt.want)
		})
	}
}

func TestEvolve(t *testing.T) {
	c := Circuit{N: 1, Gates: []Gate{{Kind: H}}}
	state := []complex128{1 + 1i, 1 - 1i}
	if err := c.Evolve(state); err != nil {
		t.Fatal(err)
	}
	checkAmplitudes(t, state, []complex128{math.Sqrt2, math.Sqrt2 * 1i})
}

func TestSimulationValidation(t *testing.T) {
	for _, n := range []int{-1, 0, MaxSimulationQubits + 1} {
		c := Circuit{N: n}
		if _, err := c.Simulate(0); err == nil {
			t.Errorf("Simulate accepted %d qubits", n)
		}
	}
	largest := Circuit{N: MaxSimulationQubits}
	last := (1 << MaxSimulationQubits) - 1
	state, err := largest.Simulate(last)
	if err != nil || len(state) != last+1 || state[last] != 1 {
		t.Fatalf("largest supported basis state failed: %v", err)
	}

	c := Circuit{N: 1, Gates: []Gate{{Kind: H}}}
	for _, basis := range []int{-1, 2} {
		if _, err := c.Simulate(basis); err == nil {
			t.Errorf("Simulate accepted basis %d", basis)
		}
	}
	if err := c.Evolve([]complex128{1}); err == nil {
		t.Error("Evolve accepted the wrong state length")
	}
	c.Gates = append(c.Gates, Gate{Kind: H, Q0: 1})
	state = []complex128{1, 0}
	if err := c.Evolve(state); err == nil {
		t.Error("Evolve accepted an invalid gate")
	}
	if !slices.Equal(state, []complex128{1, 0}) {
		t.Errorf("Evolve mutated state on error: %v", state)
	}
}

func checkAmplitudes(t *testing.T, got, want []complex128) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d amplitudes, want %d", len(got), len(want))
	}
	for i := range want {
		if diff := cmplx.Abs(got[i] - want[i]); !(diff < 1e-12) {
			t.Errorf("amplitude %d = %v, want %v", i, got[i], want[i])
		}
	}
}
