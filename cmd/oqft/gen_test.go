package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGen(t *testing.T) {
	const header = "OPENQASM 3.0;\ninclude \"stdgates.inc\";\nqubit[2] q;\n"
	const qft = header + "h q[1];\ncp(1.5707963267948966) q[0], q[1];\nh q[0];\nswap q[0], q[1];\n"
	const aqft = header + "h q[1];\nh q[0];\nswap q[0], q[1];\n"
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"-n", "2"}, qft},
		{[]string{"--type", "qft", "-n", "2"}, qft},
		{[]string{"--type", "aqft", "-n", "2"}, qft},
		{[]string{"--type", "aqft", "-n", "2", "--cutoff", "0"}, aqft},
	} {
		var out, diagnostics bytes.Buffer
		code := run(append([]string{"gen"}, tc.args...), &out, &diagnostics)
		if code != 0 || out.String() != tc.want || diagnostics.Len() != 0 {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", tc.args, code, &out, &diagnostics)
		}
	}
}

func TestGenErrors(t *testing.T) {
	for _, args := range [][]string{
		nil, {"-n", "0"}, {"-n", "x"}, {"-n", "2", "extra"},
		{"--type", "unknown", "-n", "2"}, {"--type", "aqft", "-n", "2", "--cutoff", "-1"},
	} {
		var out, diagnostics bytes.Buffer
		code := run(append([]string{"gen"}, args...), &out, &diagnostics)
		if code != 2 || out.Len() != 0 || diagnostics.Len() == 0 {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", args, code, &out, &diagnostics)
		}
	}
}

func TestGenHelp(t *testing.T) {
	var out, help bytes.Buffer
	if code := run([]string{"gen", "--help"}, &out, &help); code != 0 ||
		!strings.Contains(help.String(), "--cutoff") || out.Len() != 0 {
		t.Fatalf("exit %d, stdout %q, help %q", code, &out, &help)
	}
}
