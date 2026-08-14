package sudoku

import (
	"fmt"
	"math/bits"
)

// solveKnownUniqueByProbability builds the exact probability trace for a puzzle
// that has already been proven to have exactly one completion. In that case the
// correct digit in every empty cell has probability 100%, so repeated exact
// recounts cannot add information. One solution search plus one contiguous
// candidate buffer produces the same probability semantics with bounded memory.
func solveKnownUniqueByProbability(grid Grid, emptyCells int) (SolveResult, error) {
	solution, ok := findOneSolution(grid)
	if !ok {
		return SolveResult{}, fmt.Errorf("unique puzzle could not be completed")
	}

	result := SolveResult{
		InitialSolutions: "1",
		Steps:            make([]SolveStep, 0, emptyCells),
	}
	candidatePool := make([]CandidateAnalysis, emptyCells*Size)
	candidateOffset := 0
	current := grid

	for step := 1; ; step++ {
		idx := firstEmptyCell(current)
		if idx < 0 {
			result.FinalGrid = Compact(current)
			result.Solved = true
			return result, nil
		}

		digit := solution[idx]
		allowed := candidateMask(current, idx)
		if digit == 0 || allowed&(uint16(1)<<digit) == 0 {
			return SolveResult{}, fmt.Errorf("unique solution produced invalid digit for r%dc%d", idx/Size+1, idx%Size+1)
		}

		candidates := candidatePool[candidateOffset : candidateOffset+Size]
		candidateOffset += Size
		for candidateDigit := uint8(1); candidateDigit <= 9; candidateDigit++ {
			remaining := "0"
			probability := float64(0)
			if candidateDigit == digit {
				remaining = "1"
				probability = 100
			}
			candidates[candidateDigit-1] = CandidateAnalysis{
				Digit:              candidateDigit,
				LocallyAllowed:     allowed&(uint16(1)<<candidateDigit) != 0,
				RemainingSolutions: remaining,
				ProbabilityPercent: probability,
			}
		}

		cell := CellAnalysis{
			Row:                           idx/Size + 1,
			Column:                        idx%Size + 1,
			Candidates:                    candidates,
			RecommendedDigit:              digit,
			RecommendedRemainingSolutions: "1",
			RecommendedProbabilityPercent: 100,
			Guaranteed:                    true,
		}
		result.Steps = append(result.Steps, SolveStep{
			Step:                     step,
			RemainingSolutionsBefore: "1",
			RemainingSolutionsAfter:  "1",
			Cell:                     cell,
		})
		current[idx] = digit
	}
}

func findOneSolution(grid Grid) (Grid, bool) {
	return findOneSolutionState(newSearchState(grid))
}

func findOneSolutionState(state searchState) (Grid, bool) {
	if !propagateStateSingles(&state) {
		return Grid{}, false
	}

	bestIdx, bestMask := selectBranchCell(&state)
	if bestIdx == -1 {
		return state.grid, true
	}

	for bestMask != 0 {
		bit := bestMask & (^bestMask + 1)
		bestMask &^= bit
		digit := uint8(bits.TrailingZeros16(bit))
		next := state
		next.place(bestIdx, digit)
		if solution, ok := findOneSolutionState(next); ok {
			return solution, true
		}
	}
	return Grid{}, false
}
