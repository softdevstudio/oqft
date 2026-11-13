"""Load OQFT circuits into Qiskit and compare them with Qiskit's own QFT.

Needs: pip install qiskit qiskit_qasm3_import, and oqft on the PATH.
Prints the results recorded in test/results/qiskit-load.txt.
"""
import subprocess
import time

import numpy as np
from qiskit import qasm3, transpile
from qiskit.quantum_info import Operator
from qiskit.synthesis import synth_qft_full

BASIS = ["cx", "rz", "sx", "x"]


def gen(*args):
    return subprocess.run(["oqft", "gen", *args], check=True, capture_output=True, text=True).stdout


small = qasm3.loads(gen("-n", "6"))
same = np.allclose(Operator(small).data, Operator(synth_qft_full(6, do_swaps=True)).data)
print(f"6-qubit oqft QFT equals Qiskit's QFT matrix: {same}")

start = time.time()
big = qasm3.loads(gen("--type", "oqft", "-n", "256", "--eps", "0.01"))
print(f"256-qubit OQFT loaded in {time.time() - start:.0f} s: {dict(big.count_ops())}, depth {big.depth()}")

ours = transpile(big, basis_gates=BASIS, optimization_level=1)
print(f"OQFT after transpile: cx {ours.count_ops().get('cx')}, depth {ours.depth()}")
ref = transpile(synth_qft_full(256, do_swaps=True), basis_gates=BASIS, optimization_level=1)
print(f"Qiskit textbook QFT after transpile: cx {ref.count_ops().get('cx')}, depth {ref.depth()}")
