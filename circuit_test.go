package oqft

import (
	"math"
	"strings"
	"testing"
)

func TestNewCircuit(t *testing.T) {
	for _, n := range []int{0, -1} {
		if _, err := NewCircuit(n); err == nil {
			t.Errorf("NewCircuit(%d): no error", n)
		}
	}
	c, err := NewCircuit(3)
	if err != nil || c.N != 3 || c.Len() != 0 || c.Depth() != 0 {
		t.Fatalf("NewCircuit(3) = %+v, %v", c, err)
	}
}

func TestKindNames(t *testing.T) {
	want := map[GateKind]string{H: "h", CPhase: "cp", Swap: "swap", RZ: "rz"}
	for k, s := range want {
		if k.String() != s {
			t.Errorf("%d: got %q, want %q", int(k), k.String(), s)
		}
	}
	if len(Kinds) != len(want) {
		t.Errorf("Kinds has %d entries, want %d", len(Kinds), len(want))
	}
	if GateKind(99).Qubits() != 0 || !strings.Contains(GateKind(99).String(), "99") {
		t.Error("unknown kind not reported")
	}
}

func TestCheckRejects(t *testing.T) {
	bad := []struct {
		g    Gate
		want string
	}{
		{Gate{Kind: H, Q0: 3}, "out of range"},
		{Gate{Kind: H, Q0: -1}, "out of range"},
		{Gate{Kind: CPhase, Q0: 0, Q1: 3, Angle: 1}, "out of range"},
		{Gate{Kind: Swap, Q0: 1, Q1: 1}, "both qubits"},
		{Gate{Kind: H, Q0: 0, Q1: 2}, "second qubit"},
		{Gate{Kind: RZ, Q0: 0, Angle: math.NaN()}, "not finite"},
		{Gate{Kind: CPhase, Q0: 0, Q1: 1, Angle: math.Inf(1)}, "not finite"},
		{Gate{Kind: H, Q0: 0, Angle: 0.5}, "takes no angle"},
		{Gate{Kind: Swap, Q0: 0, Q1: 1, Angle: 0.5}, "takes no angle"},
		{Gate{Kind: GateKind(42)}, "unknown gate kind"},
	}
	c, _ := NewCircuit(3)
	for _, b := range bad {
		err := c.Add(b.g)
		if err == nil || !strings.Contains(err.Error(), b.want) {
			t.Errorf("%+v: got %v, want error containing %q", b.g, err, b.want)
		}
	}
	if c.Len() != 0 {
		t.Errorf("rejected gates were added: %d", c.Len())
	}
}

func TestHelpersPanicOnBadGate(t *testing.T) {
	c, _ := NewCircuit(2)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("no panic")
		}
	}()
	c.CP(0, 0, 1)
}

func TestCountsAndDepth(t *testing.T) {
	c, _ := NewCircuit(4)
	// Layer 1: H0, H1, H2, H3 in parallel.
	for q := 0; q < 4; q++ {
		c.H(q)
	}
	// Layer 2: CP(0,1) and CP(2,3) in parallel.
	c.CP(0, 1, math.Pi/2).CP(2, 3, math.Pi/4)
	// Layer 3: swap(1,2) waits for both pairs.
	c.Swap(1, 2)
	// Layer 3 too: RZ on qubit 0 only waits for CP(0,1).
	c.RZ(0, 0.25)
	// Layer 4: RZ on qubit 2 waits for the swap.
	c.RZ(2, -0.25)

	if c.Len() != 9 {
		t.Errorf("Len = %d, want 9", c.Len())
	}
	want := map[GateKind]int{H: 4, CPhase: 2, Swap: 1, RZ: 2}
	got := c.Counts()
	for k, n := range want {
		if got[k] != n {
			t.Errorf("Counts[%s] = %d, want %d", k, got[k], n)
		}
	}
	if c.TwoQubit() != 3 {
		t.Errorf("TwoQubit = %d, want 3", c.TwoQubit())
	}
	if c.Depth() != 4 {
		t.Errorf("Depth = %d, want 4", c.Depth())
	}
	if g := c.Gates[4]; g.Kind != CPhase || g.Q0 != 0 || g.Q1 != 1 || g.Angle != math.Pi/2 {
		t.Errorf("gate 4 = %+v", g)
	}
}

func TestDepthSerialChain(t *testing.T) {
	c, _ := NewCircuit(1)
	for i := 0; i < 5; i++ {
		c.H(0)
	}
	if c.Depth() != 5 {
		t.Errorf("Depth = %d, want 5", c.Depth())
	}
}
