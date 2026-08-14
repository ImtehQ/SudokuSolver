# SudokuSolver

SudokuSolver is a command-line Sudoku analyzer inspired by an older private project. Instead of treating a puzzle as a yes/no solvability check or using randomized attempts, it measures the **exact remaining solution space**.

For the current grid it counts how many complete valid Sudokus are still compatible with every filled cell. For the next empty cell, it then asks the same question for each possible digit.

For a candidate digit:

```text
probability = valid completions containing that digit / all remaining valid completions
```

If the current grid has 1,000 valid completions and placing `4` in the next cell leaves 816 of them, that candidate has an exact solution-space probability of 81.6%.

## A useful uniqueness detail

A standard published Sudoku is normally designed to have exactly one solution. If all of its original clues are supplied, the exact remaining solution count is therefore already `1`. In that case the correct digit for every empty cell is 100% and every other digit is 0%.

The probability distribution becomes more varied on grids that still admit multiple completions. SudokuSolver still supports those grids: it shows how the remaining completions split between candidates and marks whether the recommended digit is guaranteed or merely the largest branch.

## Usage

Provide an 81-cell puzzle using digits `1-9` for filled cells and `0` or `.` for blanks:

```bash
sudokusolver "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
```

The default analysis reports:

- current givens and empty cells;
- local candidate pressure;
- the exact number of remaining valid completions;
- whether the completion is unique;
- the exact 1-9 distribution for the next empty cell (row-major);
- the highest-probability candidate and whether it is guaranteed.

Read a formatted puzzle from a file:

```bash
sudokusolver --file puzzle.txt
```

Or pipe a puzzle through stdin:

```bash
cat puzzle.txt | sudokusolver
```

Machine-readable output:

```bash
sudokusolver --json --file puzzle.txt
```

## Probability-guided solve

Use `--solve` to repeatedly choose the candidate contained in the largest share of the remaining solution space:

```bash
sudokusolver --solve --file puzzle.txt
```

Each step reports the chosen cell, candidate probability, and how many valid completions remain after the choice. On a uniquely solvable Sudoku each step is guaranteed at 100%. On a multi-solution grid a choice below 100% is explicitly marked as a highest-probability branch rather than a logical certainty.

JSON is also available for solve mode:

```bash
sudokusolver --solve --json --file puzzle.txt
```

## Exactness and performance

Solution counts use arbitrary-precision integers, so the reported count does not overflow a fixed-width integer. The analysis is exact, not sampled.

Exact enumeration can become expensive for extremely underconstrained grids with very large solution spaces. The intended input is a normal Sudoku puzzle or a partially completed Sudoku of ordinary difficulty.

## Automatic benchmark results

The `Sudoku Benchmarks` GitHub Action automatically analyzes and fully solves four fixed, uniquely solvable project profiles: Easy, Medium, Hard, and Impossible. It records exact completion counts, solve steps, and wall-clock timings. These names are fixed project benchmark profiles rather than a universal Sudoku difficulty standard.

![Latest automatic Sudoku benchmark results](https://raw.githubusercontent.com/ImtehQ/SudokuSolver/benchmark-results/benchmark-results.svg)

### Standard challenge comparison

The same Action also runs a reproducible comparison against public solver-challenge corpora:

- **Norvig Top95:** all 95 canonical hard puzzles, pinned in [`benchmarks/norvig_top95.txt`](benchmarks/norvig_top95.txt).
- **Tdoku 17-clue corpus:** all 49,158 puzzles are attempted as a full-corpus exact-count challenge and run through pinned Tdoku; a deterministic 1,000-puzzle sample also provides repeatable SudokuSolver exact-count and probability-solve throughput.
- **Enjoy Sudoku forum hardest 1905 11+:** a deterministic 1,000-puzzle sample for exact uniqueness and probability-guided full solving.

Each measured workload has a 300-second cap. **A timeout is a benchmark result, not a workflow failure**: the generated table explicitly shows the timeout and continues with the other measurements. This keeps the full 49,158-puzzle challenge visible while still producing useful sample throughput on every run. The first v0.4.0 full-corpus attempt established the initial baseline by exceeding a 3,600-second cap.

For a fair machine-level comparison, the workflow builds **Tdoku at pinned commit `af426180dc53aef89b82868e7b3fdfcf42165654` on the same GitHub-hosted Ubuntu runner**. The Tdoku data archive is verified against Git blob `2ae6e4f8d021d2198069814c7db18bf11fcd9591` before the corpora are used. The deterministic sample seed is `20260814`.

![Latest standard Sudoku benchmark comparison](https://raw.githubusercontent.com/ImtehQ/SudokuSolver/benchmark-results/standard-benchmark-results.svg)

The comparison shows two deliberately separate measurements:

1. **Uniqueness proof throughput.** SudokuSolver counts the exact completion space; pinned Tdoku searches with a limit of two solutions. On these known-unique corpora, both must establish that no second completion exists, although their internal algorithms and amount of bookkeeping differ.
2. **SudokuSolver probability solve throughput.** This is additional work specific to this project: repeatedly calculate exact candidate distributions, choose the largest remaining branch, and continue to completion. It should not be interpreted as a like-for-like ordinary solver speed measurement.

The latest generated files live on the dedicated `benchmark-results` branch:

- `benchmark-results.json` / `benchmark-results.svg` — four project profiles;
- `standard-benchmark-results.json` / `standard-benchmark-results.svg` — standard challenge comparison, including explicit timeout states.

The workflow also uploads all four as run artifacts. Generated benchmark output never writes directly to `main`.

Benchmark corpus sources: [Peter Norvig's Top95](https://www.norvig.com/top95.txt) and the [Tdoku benchmark suite](https://github.com/t-dillon/tdoku).

## Downloads

Release-producing changes that reach `main` are published automatically using the version in [`VERSION`](VERSION). README-only changes do not create duplicate software releases. The GitHub Release contains standalone binaries for:

- Windows x86-64 and ARM64
- macOS x86-64 and Apple Silicon
- Linux x86-64 and ARM64

No Go installation is required to run a downloaded binary.

## Development

Requirements: Go 1.23 or newer.

```bash
go test ./...
go vet ./...
go build ./cmd/...
```

Benchmark tooling tests run with:

```bash
python3 -m unittest scripts.update_benchmark_readme_test scripts.run_standard_benchmarks_test
```

Formatting is checked with:

```bash
test -z "$(gofmt -l .)"
```
