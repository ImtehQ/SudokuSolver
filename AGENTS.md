# AGENTS.md

Repository-specific instructions for AI development agents working on SudokuSolver.

## Project Identity

- **Project name:** SudokuSolver
- **GitHub repository:** `ImtehQ/SudokuSolver`
- **Main/release branch:** `main`
- **Distribution environment:** GitHub Releases
- **Production server:** none; this is a standalone CLI

## Product

SudokuSolver analyzes the exact Sudoku solution space rather than using randomized completion trials. For the current grid it counts every valid complete Sudoku compatible with the filled cells. For the next empty cell (row-major), it counts how many of those remaining completions contain each digit 1-9 and reports:

```text
candidate probability = completions containing candidate / all remaining completions
```

`--solve` repeatedly chooses the candidate contained in the largest share of the remaining solution space. For a uniquely solvable Sudoku every correct choice is 100% and guaranteed. On a multi-solution grid a choice below 100% is a heuristic branch and must be described as such, never as certainty.

Counts are exact and use arbitrary-precision integers. Very underconstrained grids can be computationally expensive; do not replace exact counting with sampling without an explicit product decision.

## Production Safety

There is currently no production database, upload store, persistent volume, or server-side application state. GitHub source history and published release artifacts must still be preserved. Never rewrite `main`, delete releases/tags as a routine fix, or expose credentials.

## Application

- **Language/runtime:** Go 1.23+
- **Package manager:** Go modules
- **CLI entry point:** `cmd/sudokusolver`
- **Batch benchmark entry point:** `cmd/sudokubench`
- **Core analysis package:** `internal/sudoku`
- **Database:** none
- **Authentication:** none
- **External runtime services:** none

## Development

Install dependencies/build metadata:

```bash
go mod download
```

Run locally:

```bash
go run ./cmd/sudokusolver --help
```

Validation:

```bash
test -z "$(gofmt -l .)"
python3 -m unittest scripts.update_benchmark_readme_test scripts.run_standard_benchmarks_test
go vet ./...
go test ./...
go build ./cmd/...
```

There are no database migration or external-service integration commands because the application has no external services or database.

## Git Workflow

Every logical development change uses a new branch and PR. Never develop directly on `main`.

Preferred branch pattern:

```text
agent/short-description
```

PRs target `main`. Prefer squash merge when repository requirements allow it. Preserve unrelated changes and keep diffs focused.

## CI

- `.github/workflows/ci.yml` runs on pull requests and pushes to `main`.
- It checks `gofmt`, benchmark-tooling unit tests, `go vet`, `go test`, and builds all commands under `cmd/`.
- Do not claim CI succeeded until the workflow reaches a successful conclusion.

## Automatic Sudoku Benchmarks

- `.github/workflows/benchmarks.yml` runs after non-README pushes to `main` and can also be dispatched manually.
- Fixed quick fixtures live in `benchmarks/puzzles.json` and are named Easy, Medium, Hard, and Impossible. These are project benchmark profiles, not a universal difficulty standard.
- `benchmarks/fixtures_test.go` requires every quick benchmark puzzle to be valid, uniquely solvable, and solvable using guaranteed probability-guided steps.
- The workflow runs normal project validation and builds both `cmd/sudokusolver` and the in-process `cmd/sudokubench` batch runner.
- SudokuSolver benchmark processes use `GOMAXPROCS=1`; primary benchmark results are intentionally single-threaded.
- Quick benchmark output is `benchmark-results.json` and `benchmark-results.svg`.
- Standard comparison output is `standard-benchmark-results.json` and `standard-benchmark-results.svg`.
- All four result files are uploaded as workflow artifacts and published to the dedicated `benchmark-results` branch. The workflow must never write benchmark results directly to `main`.
- README embeds both SVGs from the generated-results branch.
- Benchmark-result branch updates must not trigger a software release or recursive benchmark run.

### Standard comparison corpus

The standard benchmark is deliberately reproducible and must not silently drift:

- **Norvig Top95:** exactly 95 puzzles committed at `benchmarks/norvig_top95.txt` from Peter Norvig's canonical Top95 set.
- **Tdoku 17-clue:** exactly 49,158 puzzles remain the source corpus, but routine automatic comparison uses a deterministic **10-puzzle sample** for both exact-count and probability-guided solve throughput.
- **Forum hardest 1905 11+:** deterministic **10-puzzle sample** from Tdoku's benchmark archive.
- **Sample seed:** `20260814`.
- **Pinned Tdoku commit:** `af426180dc53aef89b82868e7b3fdfcf42165654`.
- **Pinned Tdoku `data.zip` Git blob:** `2ae6e4f8d021d2198069814c7db18bf11fcd9591`.

The full 49,158-puzzle 17-clue challenge is optional and must only be run explicitly with `scripts/run_standard_benchmarks.py --include-full-min17`; do not make it part of every merge benchmark. The workflow verifies both the Tdoku commit and archive blob before benchmarking. Changing the Tdoku revision, corpus, sample seed, or sample size is a benchmark-definition change and requires an explicit development PR with updated documentation.

The same-runner comparison has two separate metrics and they must not be conflated:

1. SudokuSolver exact completion counting vs pinned Tdoku searching up to two solutions to prove uniqueness on known-unique corpora.
2. SudokuSolver probability-guided full solving, which performs additional exact candidate-distribution work at every step and is not equivalent to an ordinary fastest-solution benchmark.

## Releases / Deployment

- `.github/workflows/release.yml` is the normal production distribution process.
- It runs automatically on non-README pushes to `main`.
- The version comes from `VERSION`; every new release-producing change must bump this semantic version before merge.
- README-only changes are intentionally excluded from the release trigger.
- The workflow refuses to reuse an existing tag.
- It creates standalone Windows, macOS, and Linux archives plus SHA-256 checksums.
- It smoke-tests the Linux x86-64 build before creating the GitHub Release.
- Do not create a second manual deployment when the merge-triggered release workflow is already running.

A deployment is verified only after the Release workflow succeeds and the GitHub Release exists for the intended `main` commit with the expected assets.

## Persistent Data / Backups

- **Database:** none
- **Uploads:** none
- **Persistent application state:** none
- **Release artifacts:** persistent distribution artifacts; do not delete routinely
- **Source history:** Git/GitHub history is persistent project history
- **Benchmark results:** generated data on the `benchmark-results` branch; safe to regenerate from pinned corpus definitions and source

No application-data backup or restore mechanism is required because the CLI stores no production data. Do not describe source control as a tested disaster-recovery backup for user data.

## Rollback

There is no mutable server deployment to roll back. A software rollback means directing users to a previously published immutable GitHub Release or publishing a new corrected version through a new branch/PR. Never delete or rewrite current repository history as rollback behavior.

## Security-sensitive areas

Take special care with:

- GitHub Actions permissions;
- release credentials/tokens;
- benchmark workflow write access to the dedicated results branch;
- benchmark supply-chain inputs and pinned external revisions;
- parsing untrusted puzzle/file input;
- future network or update-check functionality if introduced.

Never commit secrets.

## Definition of Done

For an applicable development request:

```text
inspect -> branch -> implement -> test -> validate -> commit -> push -> PR -> CI -> merge -> release -> release verification
```

For changes affecting the solver or benchmark automation, also verify the post-merge benchmark workflow reaches a final result, the `benchmark-results` branch contains results for the intended source commit, and the README references the live result assets.

If permissions, mandatory review, Actions policy, or another external gate blocks a stage, stop at that exact stage and report what completed, what is blocked, production impact, next action, and whether a human must act.
