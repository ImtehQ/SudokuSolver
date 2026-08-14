package sudoku

import (
	"fmt"
	"hash/fnv"
	"math/bits"
	"math/rand"
)

const allDigitsMask uint16 = 0x3FE // bits 1..9

// Analysis summarizes the current constraints and the randomized completion experiment.
type Analysis struct {
	Givens                int     `json:"givens"`
	EmptyCells            int     `json:"empty_cells"`
	ForcedCells           int     `json:"forced_cells"`
	MinCandidates         int     `json:"min_candidates"`
	MaxCandidates         int     `json:"max_candidates"`
	AverageCandidates     float64 `json:"average_candidates"`
	Trials                int     `json:"trials"`
	SuccessfulCompletions int     `json:"successful_completions"`
	EstimatedProbability  float64 `json:"estimated_probability_percent"`
	Seed                  int64   `json:"seed"`
	VerifiedSolvable      *bool   `json:"verified_solvable,omitempty"`
}

// Analyze runs a reproducible randomized constraint-completion experiment.
// The result is an empirical likelihood score, not a proof that a solution exists.
func Analyze(grid Grid, trials int, seed int64, verify bool) (Analysis, error) {
	if err := Validate(grid); err != nil {
		return Analysis{}, err
	}
	if trials <= 0 {
		return Analysis{}, fmt.Errorf("trials must be greater than zero")
	}

	analysis, viable := summarize(grid)
	if seed == 0 {
		seed = defaultSeed(grid)
	}
	analysis.Seed = seed
	analysis.Trials = trials

	if !viable {
		analysis.EstimatedProbability = 0
		if verify {
			v := false
			analysis.VerifiedSolvable = &v
		}
		return analysis, nil
	}

	if analysis.EmptyCells == 0 {
		analysis.SuccessfulCompletions = trials
		analysis.EstimatedProbability = 100
	} else {
		rng := rand.New(rand.NewSource(seed))
		for i := 0; i < trials; i++ {
			if randomizedCompletion(grid, rng) {
				analysis.SuccessfulCompletions++
			}
		}
		analysis.EstimatedProbability = float64(analysis.SuccessfulCompletions) / float64(trials) * 100
	}

	if verify {
		v := HasSolution(grid)
		analysis.VerifiedSolvable = &v
	}
	return analysis, nil
}

func summarize(grid Grid) (Analysis, bool) {
	var result Analysis
	result.MinCandidates = 10
	var totalCandidates int

	for idx, value := range grid {
		if value != 0 {
			result.Givens++
			continue
		}
		result.EmptyCells++
		count := bits.OnesCount16(candidateMask(grid, idx))
		if count == 0 {
			return result, false
		}
		if count == 1 {
			result.ForcedCells++
		}
		if count < result.MinCandidates {
			result.MinCandidates = count
		}
		if count > result.MaxCandidates {
			result.MaxCandidates = count
		}
		totalCandidates += count
	}

	if result.EmptyCells == 0 {
		result.MinCandidates = 0
		result.MaxCandidates = 0
		result.AverageCandidates = 0
		return result, true
	}
	result.AverageCandidates = float64(totalCandidates) / float64(result.EmptyCells)

	return result, unitsViable(grid)
}

func candidateMask(grid Grid, idx int) uint16 {
	if grid[idx] != 0 {
		return 0
	}
	row := idx / Size
	col := idx % Size
	box := (row/3)*3 + col/3

	used := uint16(0)
	for _, cell := range rowCells(row) {
		if v := grid[cell]; v != 0 {
			used |= 1 << v
		}
	}
	for _, cell := range colCells(col) {
		if v := grid[cell]; v != 0 {
			used |= 1 << v
		}
	}
	for _, cell := range boxCells(box) {
		if v := grid[cell]; v != 0 {
			used |= 1 << v
		}
	}
	return allDigitsMask &^ used
}

func unitsViable(grid Grid) bool {
	for i := 0; i < Size; i++ {
		if !unitViable(grid, rowCells(i)) || !unitViable(grid, colCells(i)) || !unitViable(grid, boxCells(i)) {
			return false
		}
	}
	return true
}

func unitViable(grid Grid, cells [Size]int) bool {
	present := uint16(0)
	available := uint16(0)
	for _, idx := range cells {
		if value := grid[idx]; value != 0 {
			present |= 1 << value
		} else {
			available |= candidateMask(grid, idx)
		}
	}
	missing := allDigitsMask &^ present
	return missing&^available == 0
}

func randomizedCompletion(start Grid, rng *rand.Rand) bool {
	grid := start
	for {
		if !propagateSingles(&grid) {
			return false
		}

		minCount := 10
		choices := make([]int, 0, 8)
		for idx, value := range grid {
			if value != 0 {
				continue
			}
			mask := candidateMask(grid, idx)
			count := bits.OnesCount16(mask)
			if count == 0 {
				return false
			}
			if count < minCount {
				minCount = count
				choices = choices[:0]
				choices = append(choices, idx)
			} else if count == minCount {
				choices = append(choices, idx)
			}
		}

		if len(choices) == 0 {
			return Validate(grid) == nil
		}
		idx := choices[rng.Intn(len(choices))]
		mask := candidateMask(grid, idx)
		digits := digitsFromMask(mask)
		grid[idx] = digits[rng.Intn(len(digits))]
	}
}

func propagateSingles(grid *Grid) bool {
	for {
		changed := false
		for idx, value := range grid {
			if value != 0 {
				continue
			}
			mask := candidateMask(*grid, idx)
			count := bits.OnesCount16(mask)
			if count == 0 {
				return false
			}
			if count == 1 {
				grid[idx] = uint8(bits.TrailingZeros16(mask))
				changed = true
			}
		}
		if !unitsViable(*grid) {
			return false
		}
		if !changed {
			return true
		}
	}
}

func digitsFromMask(mask uint16) []uint8 {
	digits := make([]uint8, 0, 9)
	for digit := uint8(1); digit <= 9; digit++ {
		if mask&(1<<digit) != 0 {
			digits = append(digits, digit)
		}
	}
	return digits
}

func defaultSeed(grid Grid) int64 {
	h := fnv.New64a()
	_, _ = h.Write(grid[:])
	return int64(h.Sum64() & 0x7fffffffffffffff)
}

// HasSolution checks whether at least one valid completion exists without returning that completion.
func HasSolution(grid Grid) bool {
	if err := Validate(grid); err != nil || !unitsViable(grid) {
		return false
	}
	copyGrid := grid
	return search(&copyGrid)
}

func search(grid *Grid) bool {
	if !propagateSingles(grid) {
		return false
	}

	bestIdx := -1
	bestMask := uint16(0)
	bestCount := 10
	for idx, value := range grid {
		if value != 0 {
			continue
		}
		mask := candidateMask(*grid, idx)
		count := bits.OnesCount16(mask)
		if count < bestCount {
			bestIdx = idx
			bestMask = mask
			bestCount = count
		}
	}
	if bestIdx == -1 {
		return Validate(*grid) == nil
	}

	for _, digit := range digitsFromMask(bestMask) {
		next := *grid
		next[bestIdx] = digit
		if search(&next) {
			return true
		}
	}
	return false
}
