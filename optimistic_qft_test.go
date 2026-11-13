package oqft

import (
	"math"
	"math/cmplx"
	"slices"
	"testing"
)

func TestOptimisticQFTCommutation(t *testing.T) {
	for _, tc := range []struct {
		n, width int
		exact    bool
	}{
		{1, 3, true}, {5, 3, true}, {6, 2, true}, {7, 2, false}, {9, 2, false},
	} {
		baseline, err := blockedQFT(tc.n, tc.width)
		if err != nil {
			t.Fatal(err)
		}
		optimistic, err := OQFTWithBlockWidth(tc.n, tc.width)
		if err != nil {
			t.Fatal(err)
		}
		if optimistic.N != tc.n {
			t.Fatalf("got %d qubits, want %d", optimistic.N, tc.n)
		}
		want, err := baseline.Unitary()
		if err != nil {
			t.Fatal(err)
		}
		got, err := optimistic.Unitary()
		if err != nil {
			t.Fatal(err)
		}
		var squaredError float64
		for output, row := range got {
			for input, amplitude := range row {
				delta := cmplx.Abs(amplitude - want[output][input])
				squaredError += delta * delta
			}
		}
		distance := squaredError / float64(len(want))
		if tc.exact && !(distance < 1e-24) {
			t.Errorf("%+v: normalized squared error %g, want zero", tc, distance)
		}
		if !tc.exact && !(distance > 1e-6 && distance < 0.5) {
			t.Errorf("%+v: normalized squared error %g, want a small nonzero error", tc, distance)
		}
	}
}

func TestOptimisticQFTDepth(t *testing.T) {
	const width = 4
	small, err := OQFTWithBlockWidth(8*width, width)
	if err != nil {
		t.Fatal(err)
	}
	large, err := OQFTWithBlockWidth(512*width, width)
	if err != nil {
		t.Fatal(err)
	}
	if large.Depth() != small.Depth() || large.Depth() > 11*width {
		t.Fatalf("depth grew with the number of blocks or exceeded the stage bound: %d, %d", small.Depth(), large.Depth())
	}
}

func TestOptimisticQFTValidation(t *testing.T) {
	for _, input := range [][2]int{{0, 2}, {7, 0}, {7, -1}} {
		if _, err := OQFTWithBlockWidth(input[0], input[1]); err == nil {
			t.Errorf("OQFTWithBlockWidth(%d, %d) accepted invalid input", input[0], input[1])
		}
	}
}

func TestOptimisticQFTThreeSingleQubitBlocks(t *testing.T) {
	c, err := OQFTWithBlockWidth(3, 1)
	if err != nil {
		t.Fatal(err)
	}
	// The middle phase couples blocks 0 and 1. The last phase uses
	// block 2's inverse QFT to control the output on block 1.
	want := []Gate{
		{Kind: H, Q0: 2}, {Kind: H, Q0: 0},
		{Kind: CPhase, Q0: 1, Q1: 2, Angle: math.Pi / 2},
		{Kind: H, Q0: 1}, {Kind: H, Q0: 0},
		{Kind: CPhase, Q0: 0, Q1: 1, Angle: math.Pi / 2},
		{Kind: H, Q0: 0}, {Kind: Swap, Q0: 0, Q1: 2},
	}
	if len(c.Gates) != len(want) {
		t.Fatalf("got %d gates, want %d", len(c.Gates), len(want))
	}
	for i, gate := range c.Gates {
		if gate != want[i] {
			t.Errorf("gate %d: got %+v, want %+v", i, gate, want[i])
		}
	}
}

func TestOQFTError(t *testing.T) {
	for n := 1; n <= MaxUnitaryQubits; n++ {
		exact, err := QFT(n)
		if err != nil {
			t.Fatal(err)
		}
		circuit, err := OQFT(n, 0.01)
		if err != nil {
			t.Fatal(err)
		}
		for basis := 0; basis < 1<<n; basis++ {
			got, err := circuit.Simulate(basis)
			if err != nil {
				t.Fatal(err)
			}
			want, err := exact.Simulate(basis)
			if err != nil {
				t.Fatal(err)
			}
			checkAmplitudes(t, got, want)
		}
	}
	// Exercise a genuinely commuted circuit and compare to the exact QFT.
	exact, err := QFT(10)
	if err != nil {
		t.Fatal(err)
	}
	circuit, err := OQFTWithBlockWidth(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	var squaredError float64
	for basis := 0; basis < 1<<10; basis++ {
		got, err := circuit.Simulate(basis)
		if err != nil {
			t.Fatal(err)
		}
		want, err := exact.Simulate(basis)
		if err != nil {
			t.Fatal(err)
		}
		for i, value := range got {
			delta := cmplx.Abs(value - want[i])
			squaredError += delta * delta
		}
	}
	if distance := squaredError / (1 << 10); !(distance > 0 && distance < 0.25) {
		t.Fatalf("explicit-width normalized squared error = %g", distance)
	}
}

func TestOQFTSizingAndDepth(t *testing.T) {
	for _, eps := range []float64{0, -1, 2, math.NaN(), math.Inf(1)} {
		if _, err := OQFT(4, eps); err == nil {
			t.Errorf("accepted epsilon %g", eps)
		}
	}
	if _, err := OQFT(0, 0.01); err == nil {
		t.Error("accepted zero qubits")
	}
	for n := 128; n <= 2048; n *= 2 {
		width, err := OQFTBlockWidth(n, 0.01)
		if err != nil {
			t.Fatal(err)
		}
		circuit, err := OQFT(n, 0.01)
		if err != nil {
			t.Fatal(err)
		}
		if width >= n/2 || circuit.Depth() > 11*width || circuit.Depth() >= 2*n {
			t.Fatalf("%d qubits: block width %d, depth %d", n, width, circuit.Depth())
		}
	}
}

func TestOQFTSelectedWidthSampleError(t *testing.T) {
	const n, eps = 17, 1.0
	circuit, err := OQFT(n, eps)
	if err != nil {
		t.Fatal(err)
	}
	exact, err := QFT(n)
	if err != nil {
		t.Fatal(err)
	}
	if slices.Equal(circuit.Gates, exact.Gates) {
		t.Fatal("expected the commuted circuit")
	}
	inputs := []int{0, 1, 0x15555, 0x0aaaa, 0x100ff, 0x1ff00, 0x1ffff}
	var squaredError float64
	for _, input := range inputs {
		got, err := circuit.Simulate(input)
		if err != nil {
			t.Fatal(err)
		}
		want, err := exact.Simulate(input)
		if err != nil {
			t.Fatal(err)
		}
		for i, value := range got {
			delta := cmplx.Abs(value - want[i])
			squaredError += delta * delta
		}
	}
	if mean := squaredError / float64(len(inputs)); !(mean > 0 && mean < eps) {
		t.Fatalf("sample mean squared error %g, expected between zero and %g", mean, eps)
	}
}
