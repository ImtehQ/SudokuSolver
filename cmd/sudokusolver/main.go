package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ImtehQ/SudokuSolver/internal/sudoku"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("sudokusolver", flag.ContinueOnError)
	fs.SetOutput(stderr)

	filePath := fs.String("file", "", "read the puzzle from a file")
	trials := fs.Int("trials", 5000, "number of randomized completion trials")
	seed := fs.Int64("seed", 0, "random seed (0 derives a stable seed from the puzzle)")
	verify := fs.Bool("verify", false, "exactly verify whether at least one solution exists")
	jsonOutput := fs.Bool("json", false, "write analysis as JSON")
	showVersion := fs.Bool("version", false, "print version information")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "SudokuSolver analyzes Sudoku constraints and estimates completion likelihood.\n\n")
		fmt.Fprintf(stderr, "Usage:\n  sudokusolver [options] <81-cell-puzzle>\n  sudokusolver [options] --file puzzle.txt\n  cat puzzle.txt | sudokusolver [options]\n\nOptions:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		fmt.Fprintf(stdout, "SudokuSolver %s (%s)\n", version, commit)
		return 0
	}
	if *trials < 1 || *trials > 1_000_000 {
		fmt.Fprintln(stderr, "error: --trials must be between 1 and 1000000")
		return 2
	}

	input, err := readPuzzleInput(fs.Args(), *filePath, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 2
	}
	grid, err := sudoku.Parse(input)
	if err != nil {
		fmt.Fprintf(stderr, "invalid puzzle: %v\n", err)
		return 2
	}
	analysis, err := sudoku.Analyze(grid, *trials, *seed, *verify)
	if err != nil {
		fmt.Fprintf(stderr, "analysis failed: %v\n", err)
		return 1
	}

	if *jsonOutput {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(analysis); err != nil {
			fmt.Fprintf(stderr, "output failed: %v\n", err)
			return 1
		}
		return 0
	}

	fmt.Fprintf(stdout, "SudokuSolver analysis\n")
	fmt.Fprintf(stdout, "Givens: %d\n", analysis.Givens)
	fmt.Fprintf(stdout, "Empty cells: %d\n", analysis.EmptyCells)
	fmt.Fprintf(stdout, "Forced cells now: %d\n", analysis.ForcedCells)
	fmt.Fprintf(stdout, "Candidate range: %d-%d\n", analysis.MinCandidates, analysis.MaxCandidates)
	fmt.Fprintf(stdout, "Average candidates: %.2f\n", analysis.AverageCandidates)
	fmt.Fprintf(stdout, "Randomized trials: %d\n", analysis.Trials)
	fmt.Fprintf(stdout, "Successful completions: %d\n", analysis.SuccessfulCompletions)
	fmt.Fprintf(stdout, "Estimated solution probability: %.2f%%\n", analysis.EstimatedProbability)
	if analysis.VerifiedSolvable != nil {
		if *analysis.VerifiedSolvable {
			fmt.Fprintln(stdout, "Exact solvability: yes")
		} else {
			fmt.Fprintln(stdout, "Exact solvability: no")
		}
	}
	fmt.Fprintln(stdout, "Note: the percentage is an empirical randomized completion score, not a mathematical proof of solvability. Use --verify for an exact yes/no check.")
	return 0
}

func readPuzzleInput(args []string, filePath string, stdin io.Reader) (string, error) {
	if filePath != "" && len(args) > 0 {
		return "", errors.New("use either --file or a positional puzzle, not both")
	}
	if len(args) > 1 {
		return "", errors.New("expected at most one positional puzzle")
	}
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("read %s: %w", filePath, err)
		}
		return string(data), nil
	}
	if len(args) == 1 {
		return strings.TrimSpace(args[0]), nil
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	if strings.TrimSpace(string(data)) == "" {
		return "", errors.New("no puzzle provided")
	}
	return string(data), nil
}
