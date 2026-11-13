package oqft

import (
	"math"
	"testing"
)

func TestCheckerStateError(t *testing.T) {
	identity, err := NewCircuit(1)
	if err != nil {
		t.Fatal(err)
	}
	checker, err := NewChecker(identity)
	if err != nil {
		t.Fatal(err)
	}
	identity.H(0)
	for basis, want := range []float64{2 - math.Sqrt2, 2 + math.Sqrt2} {
		got, err := checker.StateError(basis)
		if err != nil || math.Abs(got-want) > 1e-12 {
			t.Fatalf("basis %d: got %g, %v, want %g", basis, got, err, want)
		}
	}
	for _, basis := range []int{-1, 2} {
		if _, err := checker.StateError(basis); err == nil {
			t.Errorf("accepted basis %d", basis)
		}
	}
}

func TestCheckerValidation(t *testing.T) {
	for _, circuit := range []*Circuit{nil, {N: MaxSimulationQubits + 1}, {N: 1, Gates: []Gate{{Kind: H, Q0: 1}}}} {
		if _, err := NewChecker(circuit); err == nil {
			t.Errorf("accepted circuit %+v", circuit)
		}
	}
}

func TestCheckerWraparound(t *testing.T) {
	circuit, err := OQFTWithBlockWidth(12, 3)
	if err != nil {
		t.Fatal(err)
	}
	checker, err := NewChecker(circuit)
	if err != nil {
		t.Fatal(err)
	}
	// Block 2 is estimated using the phase contributed by block 3.
	// Hold the other blocks fixed and sweep the estimated block value.
	for value := 0; value < 8; value++ {
		basis := (value << 3) | 7
		distance, err := checker.StateError(basis)
		if err != nil {
			t.Fatal(err)
		}
		if failed, want := distance > 0.5, value == 7; failed != want {
			t.Errorf("block value %d: squared error %g, want failure %t", value, distance, want)
		}
	}
}
