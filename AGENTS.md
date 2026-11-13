# Instructions for coding agents (OpenAI Codex, ChatGPT and others)

Read CLAUDE.md first. It is the single source of truth for this repo's
workflow, plan, state and to-do list, and the same rules apply to every
agent. Claude and ChatGPT/Codex both work on this repo and hand work to
each other through CLAUDE.md alone.

The rules that matter most:

- `git pull` before you start: other agents push to main too.
- Work on main. No side branches, no pull requests to yourself.
- Before writing code, push an "In progress" section to CLAUDE.md (what
  and why, design decisions, steps done and remaining, where to resume).
  Update and push it after each finished step, so the other agent can
  continue from the repo alone. CLAUDE.md-only commits carry `[skip ci]`.
- Never create or move release tags by hand. A release happens when
  VERSION, `const Version` in version.go and a new CHANGELOG.md section are
  raised together and CI goes green on main.
- Run `gofmt -l .`, `go vet ./...` and `go test -race ./...` before pushing.
- Go standard library only (Claude's sandbox can't reach the module proxy).
- Numbers on oqft.com come from test/results, written by
  scripts/record-results.sh. Don't type numbers by hand.
- Stuck: read the real failing output before any fix or rerun. At most two
  fix rounds on one failure, then write the blocker and the options into
  CLAUDE.md and stop.
- When you finish, update the State, In progress and Gaps sections of
  CLAUDE.md with what you did and what is still open.
