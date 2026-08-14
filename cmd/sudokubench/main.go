package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/ImtehQ/SudokuSolver/internal/sudoku"
)

type benchmarkResult struct {
	Dataset                   string  `json:"dataset"`
	Mode                      string  `json:"mode"`
	Puzzles                   int     `json:"puzzles"`
	Successful                int     `json:"successful"`
	Unique                    int     `json:"unique"`
	TotalSeconds              float64 `json:"total_seconds"`
	PuzzlesPerSecond          float64 `json:"puzzles_per_second"`
	MeanMilliseconds          float64 `json:"mean_milliseconds"`
	MedianMilliseconds        float64 `json:"median_milliseconds"`
	P95Milliseconds           float64 `json:"p95_milliseconds"`
	MaxMilliseconds           float64 `json:"max_milliseconds"`
	AllocatedBytes            uint64  `json:"allocated_bytes"`
	Allocations               uint64  `json:"allocations"`
	GCCycles                  uint32  `json:"gc_cycles"`
	AllocatedBytesPerPuzzle   float64 `json:"allocated_bytes_per_puzzle"`
	AllocationsPerPuzzle      float64 `json:"allocations_per_puzzle"`
	TotalSteps                int     `json:"total_steps,omitempty"`
	AverageSteps              float64 `json:"average_steps,omitempty"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("sudokubench", flag.ContinueOnError)
	fs.SetOutput(stderr)

	inputPath := fs.String("input", "", "newline-delimited Sudoku dataset")
	dataset := fs.String("dataset", "dataset", "dataset label")
	mode := fs.String("mode", "count", "benchmark mode: count or solve")
	limit := fs.Int("limit", 0, "maximum puzzles to process (0 means all)")
	requireUnique := fs.Bool("require-unique", true, "fail if a puzzle does not have exactly one completion")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *inputPath == "" {
		fmt.Fprintln(stderr, "error: --input is required")
		return 2
	}
	if *mode != "count" && *mode != "solve" {
		fmt.Fprintln(stderr, "error: --mode must be count or solve")
		return 2
	}
	if *limit < 0 {
		fmt.Fprintln(stderr, "error: --limit must be non-negative")
		return 2
	}

	file, err := os.Open(*inputPath)
	if err != nil {
		fmt.Fprintf(stderr, "error: open dataset: %v\n", err)
		return 2
	}
	defer file.Close()

	puzzles, err := readPuzzles(file, *limit)
	if err != nil {
		fmt.Fprintf(stderr, "error: read dataset: %v\n", err)
		return 2
	}
	if len(puzzles) == 0 {
		fmt.Fprintln(stderr, "error: dataset contains no puzzles")
		return 2
	}

	result, err := benchmark(*dataset, *mode, puzzles, *requireUnique)
	if err != nil {
		fmt.Fprintf(stderr, "benchmark failed: %v\n", err)
		return 1
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(stderr, "output failed: %v\n", err)
		return 1
	}
	return 0
}

func readPuzzles(r io.Reader, limit int) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	puzzles := make([]string, 0)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		puzzle, ok, err := normalizePuzzleLine(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if !ok {
			continue
		}
		puzzles = append(puzzles, puzzle)
		if limit > 0 && len(puzzles) >= limit {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return puzzles, nil
}

func normalizePuzzleLine(line string) (string, bool, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", false, nil
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", false, nil
	}
	token := fields[0]
	if len(token) < 81 {
		return "", false, fmt.Errorf("puzzle token has %d cells, want at least 81", len(token))
	}
	token = token[:81]
	var builder strings.Builder
	builder.Grow(81)
	for _, char := range token {
		switch {
		case char == '.' || char == '0':
			builder.WriteByte('0')
		case char >= '1' && char <= '9':
			builder.WriteRune(char)
		default:
			return "", false, fmt.Errorf("invalid puzzle character %q", char)
		}
	}
	return builder.String(), true, nil
}

func benchmark(dataset, mode string, puzzles []string, requireUnique bool) (benchmarkResult, error) {
	result := benchmarkResult{Dataset: dataset, Mode: mode, Puzzles: len(puzzles)}
	durations := make([]time.Duration, len(puzzles))
	grids := make([]sudoku.Grid, len(puzzles))
	startedSetup := time.Now()
	for index, text := range puzzles {
		grid, err := sudoku.Parse(text)
		if err != nil {
			return benchmarkResult{}, fmt.Errorf("puzzle %d: parse: %w", index+1, err)
		}
		grids[index] = grid
	}
	setupDuration := time.Since(startedSetup)

	var memoryBefore runtime.MemStats
	runtime.ReadMemStats(&memoryBefore)
	startedOperations := time.Now()
	one := big.NewInt(1)

	for index, grid := range grids {
		started := time.Now()
		switch mode {
		case "count":
			count, err := sudoku.CountSolutions(grid)
			durations[index] = time.Since(started)
			if err != nil {
				return benchmarkResult{}, fmt.Errorf("puzzle %d: count: %w", index+1, err)
			}
			if count.Cmp(one) == 0 {
				result.Unique++
			}
			if requireUnique && count.Cmp(one) != 0 {
				return benchmarkResult{}, fmt.Errorf("puzzle %d: expected exactly one completion, got %s", index+1, count.String())
			}
			result.Successful++
		case "solve":
			solved, err := sudoku.SolveByProbability(grid)
			durations[index] = time.Since(started)
			if err != nil {
				return benchmarkResult{}, fmt.Errorf("puzzle %d: solve: %w", index+1, err)
			}
			if solved.InitialSolutions == "1" {
				result.Unique++
			}
			if requireUnique && solved.InitialSolutions != "1" {
				return benchmarkResult{}, fmt.Errorf("puzzle %d: expected exactly one completion, got %s", index+1, solved.InitialSolutions)
			}
			if !solved.Solved {
				return benchmarkResult{}, fmt.Errorf("puzzle %d: probability-guided solve did not finish", index+1)
			}
			for _, step := range solved.Steps {
				if requireUnique && !step.Cell.Guaranteed {
					return benchmarkResult{}, fmt.Errorf("puzzle %d: used a non-guaranteed step", index+1)
				}
			}
			result.TotalSteps += len(solved.Steps)
			result.Successful++
		}
	}
	operationDuration := time.Since(startedOperations)

	var memoryAfter runtime.MemStats
	runtime.ReadMemStats(&memoryAfter)
	result.AllocatedBytes = memoryAfter.TotalAlloc - memoryBefore.TotalAlloc
	result.Allocations = memoryAfter.Mallocs - memoryBefore.Mallocs
	result.GCCycles = memoryAfter.NumGC - memoryBefore.NumGC
	result.TotalSeconds = (setupDuration + operationDuration).Seconds()
	if result.TotalSeconds > 0 {
		result.PuzzlesPerSecond = float64(result.Puzzles) / result.TotalSeconds
	}
	if result.Puzzles > 0 {
		result.AllocatedBytesPerPuzzle = float64(result.AllocatedBytes) / float64(result.Puzzles)
		result.AllocationsPerPuzzle = float64(result.Allocations) / float64(result.Puzzles)
	}
	fillLatencyStats(&result, durations)
	if result.Puzzles > 0 && mode == "solve" {
		result.AverageSteps = float64(result.TotalSteps) / float64(result.Puzzles)
	}
	return result, nil
}

func fillLatencyStats(result *benchmarkResult, durations []time.Duration) {
	if len(durations) == 0 {
		return
	}
	values := make([]float64, len(durations))
	var total float64
	for i, duration := range durations {
		milliseconds := float64(duration) / float64(time.Millisecond)
		values[i] = milliseconds
		total += milliseconds
	}
	sort.Float64s(values)
	result.MeanMilliseconds = total / float64(len(values))
	result.MedianMilliseconds = percentile(values, 0.50)
	result.P95Milliseconds = percentile(values, 0.95)
	result.MaxMilliseconds = values[len(values)-1]
}

func percentile(sortedValues []float64, percentile float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}
	if percentile <= 0 {
		return sortedValues[0]
	}
	if percentile >= 1 {
		return sortedValues[len(sortedValues)-1]
	}
	position := percentile * float64(len(sortedValues)-1)
	lower := int(position)
	upper := lower + 1
	if upper >= len(sortedValues) {
		return sortedValues[lower]
	}
	fraction := position - float64(lower)
	return sortedValues[lower] + (sortedValues[upper]-sortedValues[lower])*fraction
}
