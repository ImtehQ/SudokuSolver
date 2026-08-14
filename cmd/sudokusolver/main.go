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
	solve := fs.Bool("solve", false, "solve by repeatedly choosing the highest exact candidate probability")
	jsonOutput := fs.Bool("json", false, "write analysis as JSON")
	showVersion := fs.Bool("version", false, "print version information")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "SudokuSolver counts the exact remaining Sudoku solution space and candidate probabilities.\n\n")
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

	if *solve {
		result, err := sudoku.SolveByProbability(grid)
		if err != nil {
			fmt.Fprintf(stderr, "solve failed: %v\n", err)
			return 1
		}
		if *jsonOutput {
			if err := writeJSON(stdout, result); err != nil {
				fmt.Fprintf(stderr, "output failed: %v\n", err)
				return 1
			}
			if !result.Solved {
				return 1
			}
			return 0
		}
		writeSolveText(stdout, result)
		if !result.Solved {
			return 1
		}
		return 0
	}

	analysis, err := sudoku.Analyze(grid)
	if err != nil {
		fmt.Fprintf(stderr, "analysis failed: %v\n", err)
		return 1
	}
	if *jsonOutput {
		if err := writeJSON(stdout, analysis); err != nil {
			fmt.Fprintf(stderr, "output failed: %v\n", err)
			return 1
		}
		return 0
	}

	writeAnalysisText(stdout, analysis)
	return 0
}

func writeJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func writeAnalysisText(w io.Writer, analysis sudoku.Analysis) {
	fmt.Fprintln(w, "SudokuSolver exact solution-space analysis")
	fmt.Fprintf(w, "Givens: %d\n", analysis.Givens)
	fmt.Fprintf(w, "Empty cells: %d\n", analysis.EmptyCells)
	fmt.Fprintf(w, "Forced cells now: %d\n", analysis.ForcedCells)
	fmt.Fprintf(w, "Candidate range: %d-%d\n", analysis.MinCandidates, analysis.MaxCandidates)
	fmt.Fprintf(w, "Average candidates: %.2f\n", analysis.AverageCandidates)
	fmt.Fprintf(w, "Remaining valid completions: %s\n", analysis.RemainingSolutions)
	fmt.Fprintf(w, "Unique completion: %s\n", yesNo(analysis.UniqueSolution))

	if analysis.RemainingSolutions == "0" {
		fmt.Fprintln(w, "No valid Sudoku completion exists for the current grid.")
		return
	}
	if analysis.NextCell == nil {
		fmt.Fprintln(w, "The grid is already complete.")
		return
	}

	cell := analysis.NextCell
	fmt.Fprintf(w, "\nNext cell: r%dc%d\n", cell.Row, cell.Column)
	fmt.Fprintln(w, "Candidate distribution:")
	for _, candidate := range cell.Candidates {
		blocked := ""
		if !candidate.LocallyAllowed {
			blocked = " [blocked by current row/column/box]"
		}
		fmt.Fprintf(w, "  %d: %s / %s (%.6f%%)%s\n",
			candidate.Digit,
			candidate.RemainingSolutions,
			analysis.RemainingSolutions,
			candidate.ProbabilityPercent,
			blocked,
		)
	}
	fmt.Fprintf(w, "Recommended digit: %d (%.6f%%; leaves %s completion(s))\n",
		cell.RecommendedDigit,
		cell.RecommendedProbabilityPercent,
		cell.RecommendedRemainingSolutions,
	)
	if cell.Guaranteed {
		fmt.Fprintln(w, "Recommendation status: guaranteed across every remaining completion.")
	} else {
		fmt.Fprintln(w, "Recommendation status: highest-probability branch, not guaranteed.")
	}
	if analysis.UniqueSolution {
		fmt.Fprintln(w, "Note: this grid already has exactly one valid completion, so its correct empty digits are 100% under exact solution-space counting.")
	}
}

func writeSolveText(w io.Writer, result sudoku.SolveResult) {
	fmt.Fprintln(w, "SudokuSolver probability-guided solve")
	fmt.Fprintf(w, "Initial remaining valid completions: %s\n", result.InitialSolutions)
	if result.InitialSolutions == "0" {
		fmt.Fprintln(w, "No valid Sudoku completion exists for the current grid.")
		return
	}

	for _, step := range result.Steps {
		status := "highest probability"
		if step.Cell.Guaranteed {
			status = "guaranteed"
		}
		fmt.Fprintf(w, "Step %d: r%dc%d = %d (%.6f%%, %s -> %s completion(s), %s)\n",
			step.Step,
			step.Cell.Row,
			step.Cell.Column,
			step.Cell.RecommendedDigit,
			step.Cell.RecommendedProbabilityPercent,
			step.RemainingSolutionsBefore,
			step.RemainingSolutionsAfter,
			status,
		)
	}

	if result.Solved {
		fmt.Fprintln(w, "\nFinal grid:")
		writeGrid(w, result.FinalGrid)
	} else {
		fmt.Fprintln(w, "The grid could not be completed.")
	}
}

func writeGrid(w io.Writer, compact string) {
	for row := 0; row < sudoku.Size; row++ {
		start := row * sudoku.Size
		fmt.Fprintln(w, compact[start:start+sudoku.Size])
	}
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
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
