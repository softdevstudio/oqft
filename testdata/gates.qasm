OPENQASM 3.0;
include "stdgates.inc";
qubit[3] q;
h q[2];
cp(0.25) q[0], q[2];
swap q[2], q[1];
rz(-0.5) q[1];
