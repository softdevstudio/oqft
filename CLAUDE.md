# OQFT: where the next session starts

This file is the single source of truth for every agent on this repo:
Claude, ChatGPT and Codex. Read it all before working. AGENTS.md points
here. There is no separate roadmap file: the plan, the state and the to-do
list live below.

OQFT is an open source toolkit for the optimistic quantum Fourier
transform (arXiv:2505.00701), plus a Q-Day calculator on oqft.com that
runs the same estimator in the browser (compiled to WebAssembly). Go,
standard library only, one binary.

- Code: github.com/softdevstudio/oqft (this repo), Apache 2.0
- Site: github.com/photozeitgeist/oqft.com (Hugo, Launch template, content
  and config values only; never touch the theme or layouts). Not created
  yet; the owner creates it before Phase 4.
- Owner: YK. His part is approvals and creating repos. He never publishes
  releases and doesn't run builds.

Layout (grows phase by phase):

- `version.go` holds `const Version`; `VERSION` and `CHANGELOG.md` must
  agree (scripts/check-version.sh, run in CI).
- `cmd/oqft` is the CLI. Library code lives in the root package `oqft`
  until it is big enough to split into `internal/` packages.
- `scripts/record-results.sh` writes `test/results/`.

## Handoff rules (Claude and ChatGPT/Codex)

1. `git pull` before working. Both agents push to main.
2. Read this file: "In progress" says where to resume.
3. Run `go test ./...`. If it fails before you've changed anything, that's
   the first job: fix it or write it up under In progress.
4. Before writing code for a new step, push an "In progress" section here:
   what and why, design decisions, steps with [x]/[ ], where to resume.
   Update and push it after each finished step with the code. Commits that
   only touch CLAUDE.md carry `[skip ci]`.
5. When a phase is done: move its In progress section to "Done", update
   State and Gaps, and say what's next. The next phase waits for the
   owner's go.
6. If blocked: write the blocker and the options under In progress, push,
   and stop. The other agent or the owner picks it up from there.
7. Work in short assignments, to keep each session within usage limits.
   When a phase starts, split it into assignments (1a, 1b, ...) in the
   plan below before any code. One assignment per session: finish it,
   push code with green CI, tick it here, and stop. The owner starts the
   next one ("next assignment"). Don't run on into the next assignment.

## Workflow (low maintenance)

- Releases: raise `VERSION`, `version.go` and add a `CHANGELOG.md` section
  in one commit, push to main. When CI is green, `release.yml` builds
  Linux, macOS and Windows archives, smoke-tests them, tags `vX.Y.Z` and
  publishes "OQFT vX.Y.Z" with SHA256SUMS. 0.0.0 is never released. Never
  tag by hand. After a release: watch the run, check the files and links,
  run a downloaded binary.
- The release runs after any CI run on main completes and builds the
  commit that run tested. Raise VERSION only in a commit you want released.
  Don't push to main while a release run is going.
- Tests: a test for every operation, plus failure tests aimed at the main
  promise (the optimistic circuit matches the true QFT except on a small,
  predicted slice of inputs). Every number on the site comes from
  `test/results/`. Re-record after any change that could move a number.
- An independent review agent looks at code and site before a phase is
  called done; then fixes and a fresh recording.
- `gofmt -l .`, `go vet ./...`, `go test -race ./...` before every push.
- Go standard library only: Claude's sandbox can't reach the Go module
  proxy. go.mod stays at go 1.24 (the sandbox has 1.24.7); release builds
  use current stable Go.
- GitHub job logs can't be downloaded from Claude's sessions, but
  annotations can, so workflows print `::error::` lines. Read the real
  error before any rerun (a blind rerun counts as a fix round). At most
  two fix rounds on one failure, or 15 minutes of diagnosis with no known
  cause: stop, write what's known and the options here, hand over.
- Site check after a deploy: once. Cloudflare Pages status on the site
  commit plus one page not fetched before. Never re-fetch the same address.
- main is the default branch. No side branches.
- Content rules for the site: the owner's Hugo post rules and plain
  writing (no AI tells). Nothing about domain names or their value in any
  project material.

## Plan

From the owner's dev plan of 2026-11-13; he started the project by creating
this repo. Each phase is split into short assignments (handoff rule 7);
each ends with green CI and this file
updated. Ask the owner for his go before starting each phase.

0. [x] Repo setup: module, `oqft version`, CI, release workflow, this file
   and AGENTS.md. No release.
1. [ ] Circuit core and simulator (0.1.0), in assignments:
   - [x] 1a Circuit model: gates H, controlled phase, swap, Z rotation;
     circuit type; gate count and depth. Done when unit tests pass.
   - [ ] 1b Statevector simulator (n <= 20). Done when single-gate and
     small-circuit tests pass.
   - [ ] 1c Textbook QFT generator and full-unitary builder (n <= 10).
     Done when the QFT equals the DFT matrix for n = 1..10.
   - [ ] 1d Approximate QFT generator (drops rotations below a cutoff).
     Done when error shrinks as the cutoff rises, checked against the
     exact QFT.
   - [ ] 1e OpenQASM 3 export; `oqft gen --type qft|aqft -n N`. Done when
     golden-file and CLI tests pass.
   - [ ] 1f `oqft stats`; scripts/record-results.sh run; README and
     CHANGELOG. Done when results are recorded in test/results.
   - [ ] 1g Independent review and fixes; raise to 0.1.0. Done when the
     release workflow publishes v0.1.0 and a downloaded binary is checked.
2. [ ] The optimistic QFT (0.2.0): first check whether the paper's authors
   published code to test against (record it in NOTICE if used). Build the
   OQFT from arXiv:2505.00701: block size m from n and eps, intra-block
   QFTs, block-phased rotations, the inserted phase-estimation blocks.
   `oqft gen --type oqft -n N --eps E`. Tests: average error (the paper's
   Frobenius-norm definition) below eps for every n the simulator handles;
   depth grows as log(n/eps), from recorded counts up to n = 2048. Hardest
   phase: the stuck rules apply strictly.
3. [ ] The checker (0.3.0): `oqft check -n N --eps E`, per-state error,
   failing states, failure map (text and CSV) showing the wraparound edges;
   random-basis check; failure tests that fail where predicted and only
   there.
4. [ ] Site launch: oqft.com on the Launch template. Pages About (incl. the
   name), Contact (info@oqft.com), Quick Start, Download, Reference. Posts:
   launch, plain-language explainer of optimistic circuits, how it was
   tested, the name. Images in static/images. Footer "A SoftDevStudio.com
   project"; OQFT added to softdevstudio.com.
5. [ ] Resource estimator (0.4.0): `oqft estimate --rsa 2048`, logical
   qubits, depth, gate counts for the QFT and arithmetic parts of Shor's
   algorithm, textbook QFT vs OQFT; physical-qubit range from a small
   surface-code model with assumptions as flags; `data/estimates.yaml`
   reference points from published papers, each with its source, and tests
   that reproduce them within a stated tolerance. No market-size figures.
6. [ ] Q-Day calculator (0.5.0): estimator built to WebAssembly; a
   standalone page in the site's `static/calculator/` (no template change,
   one menu entry); the release attaches the .wasm; a CI test checks the
   calculator and `oqft estimate` agree.
7. [ ] 1.0.0: reference docs, performance pass for n = 4096, full fresh
   recording, independent review, site numbers refreshed.

Later, not scheduled: a thin Python package for Qiskit and Cirq users;
hardware layout planner (arXiv:2605.15297); other optimistic circuits
(adders, multipliers).

## In progress

Nothing. Next is assignment 1b (statevector simulator), once the owner
says "next assignment". Start by writing its In progress section here.

## Done

### Assignment 1a: circuit model (2026-11-13, Claude)

`circuit.go`: `GateKind` (`H`, `CPhase`, `Swap`, `RZ`; `String()` gives
the OpenQASM 3 names h, cp, swap, rz; `Kinds` lists them in order),
`Gate{Kind, Q0, Q1, Angle}` with `Check(n)`, `Circuit{N, Gates}` with
`NewCircuit`, `Add` (validates, returns error), chainable helpers `H`,
`CP(control, target, theta)`, `Swap`, `RZ` (panic on a bad gate),
`Len`, `Counts`, `TwoQubit`, `Depth` (ASAP layers, one layer per gate).
Conventions for 1b: CPhase is diag(1,1,1,e^{iθ}); RZ is
diag(e^{-iθ/2}, e^{iθ/2}); qubit 0 is the least significant bit of the
basis-state index (to be fixed in 1b's simulator and kept everywhere).
Tests in `circuit_test.go`: validation errors, panics, counts, depth.

### Phase 0: repo setup (2026-11-13, Claude)

Go module github.com/softdevstudio/oqft, `oqft version` and `help`, CLI
tests, CI (version check, gofmt, vet, race tests on Linux, macOS,
Windows), release workflow copied from PXTA (skips 0.0.0), LICENSE,
NOTICE (credits the paper, no affiliation), THIRD_PARTY_LICENSES.md,
README, CHANGELOG, AGENTS.md, this file, scripts/check-version.sh and
scripts/record-results.sh.

## State

0.0.0. Circuit model in place (1a). Nothing released.

## Gaps and to-do

- Site repo photozeitgeist/oqft.com doesn't exist yet (needed for Phase 4).
