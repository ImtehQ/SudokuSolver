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

The `Sudoku Benchmarks` GitHub Action automatically analyzes and fully solves four fixed, uniquely solvable profiles: Easy, Medium, Hard, and Impossible. It records exact completion counts, solve steps, and wall-clock timings, then updates only the generated section below through its own pull request.

<!-- benchmark-results:start -->
These results are generated automatically by the `Sudoku Benchmarks` GitHub Action.
Each fixture must have exactly one valid completion, and every solve step must therefore be a guaranteed 100% choice.
The difficulty names are fixed benchmark profiles, not a universal Sudoku rating standard.
Timings come from a GitHub-hosted Ubuntu runner and will vary between runs.

Benchmark source commit: `c66e7e4e3628`

| Difficulty | Givens | Exact completions | Analysis time | Solve steps | Full solve time | Result |
|---|---:|---:|---:|---:|---:|---|
| Easy | 47 | 1 | 2 ms | 34 | 2 ms | ✅ solved |
| Medium | 30 | 1 | 1 ms | 51 | 3 ms | ✅ solved |
| Hard | 28 | 1 | 2 ms | 53 | 5 ms | ✅ solved |
| Impossible | 21 | 1 | 163 ms | 60 | 212 ms | ✅ solved |
<!-- benchmark-results:end -->

## Downloads

Release-producing changes that reach `main` are published automatically using the version in [`VERSION`](VERSION). README-only benchmark-result updates do not create duplicate software releases. The GitHub Release contains standalone binaries for:

- Windows x86-64 and ARM64
- macOS x86-64 and Apple Silicon
- Linux x86-64 and ARM64

No Go installation is required to run a downloaded binary.

## Development

Requirements: Go 1.23 or newer.

```bash
go test ./...
go vet ./...
go build ./cmd/sudokusolver
```

Benchmark tooling tests run with:

```bash
python3 -m unittest scripts.update_benchmark_readme_test
```

Formatting is checked with:

```bash
test -z "$(gofmt -l .)"
```
