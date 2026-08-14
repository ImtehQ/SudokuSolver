# SudokuSolver

SudokuSolver is a command-line Sudoku analyzer inspired by an older private project. Instead of simply printing a solved grid, it examines the current constraints and estimates how likely a randomized, constraint-respecting completion is to reach a valid Sudoku.

The project is intentionally analytical: it reports givens, empty cells, forced cells, candidate pressure, randomized completion success, and—when requested—an exact yes/no solvability verification without revealing the solution.

## Important meaning of “probability”

A fixed Sudoku position is ultimately either solvable or unsolvable, so there is no literal unknown probability once the full state is known. SudokuSolver therefore reports an **empirical solution probability**: the percentage of randomized constraint-completion trials that successfully reach a valid completed grid.

That percentage is useful as a comparative analytical score, but it is not a mathematical proof that a solution exists. Use `--verify` when you want an exact existence check.

## Usage

Provide an 81-cell puzzle using digits `1-9` for givens and `0` or `.` for blanks:

```bash
sudokusolver "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
```

Read a formatted puzzle from a file:

```bash
sudokusolver --file puzzle.txt
```

Or pipe a puzzle through stdin:

```bash
cat puzzle.txt | sudokusolver
```

Useful options:

```text
--trials N   randomized completion trials (default 5000)
--seed N     reproducible random seed; 0 derives one from the puzzle
--verify     exactly verify whether at least one solution exists
--json       emit machine-readable JSON
--version    show build version and commit
```

Example:

```bash
sudokusolver --verify --trials 10000 --file puzzle.txt
```

## Downloads

Each accepted change to `main` is released automatically using the version in [`VERSION`](VERSION). The GitHub Release contains standalone binaries for:

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

Formatting is checked with:

```bash
test -z "$(gofmt -l .)"
```
