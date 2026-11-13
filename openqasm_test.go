package oqft

import (
	"os"
	"testing"
)

func TestOpenQASM(t *testing.T) {
	c := Circuit{N: 3}
	c.H(2).CP(0, 2, 0.25).Swap(2, 1).RZ(1, -0.5)
	got, err := c.OpenQASM()
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/gates.qasm")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestOpenQASMInvalidCircuit(t *testing.T) {
	for _, c := range []Circuit{
		{},
		{N: 1, Gates: []Gate{{Kind: H}, {Kind: Swap, Q1: 1}}},
	} {
		if program, err := c.OpenQASM(); err == nil || program != "" {
			t.Errorf("OpenQASM = %q, %v; want an error and no program", program, err)
		}
	}
}
