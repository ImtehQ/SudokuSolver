package sudoku

import (
	"fmt"
	"math/big"
	"math/bits"
)

const allDigitsMask uint16 = 0x3FE // bits 1..9

// CandidateAnalysis describes how many complete valid Sudokus remain when a
// digit is placed in the analyzed cell.
type CandidateAnalysis struct {
	Digit              uint8   `json:"digit"`
	LocallyAllowed     bool    `json:"locally_allowed"`
	RemainingSolutions string  `json:"remaining_solutions"`
	ProbabilityPercent float64 `json:"probability_percent"`
}

// CellAnalysis describes the exact solution distribution for one empty cell.
type CellAnalysis struct {
	Row                           int                 `json:"row"`
	Column                        int                 `json:"column"`
	Candidates                    []CandidateAnalysis `json:"candidates"`
	RecommendedDigit              uint8               `json:"recommended_digit"`
	RecommendedRemainingSolutions string              `json:"recommended_remaining_solutions"`
	RecommendedProbabilityPercent float64             `json:"recommended_probability_percent"`
	Guaranteed                    bool                `json:"guaranteed"`
}

// Analysis summarizes the current constraints and exact remaining solution space.
type Analysis struct {
	Givens             int           `json:"givens"`
	EmptyCells         int           `json:"empty_cells"`
	ForcedCells        int           `json:"forced_cells"`
	MinCandidates      int           `json:"min_candidates"`
	MaxCandidates      int           `json:"max_candidates"`
	AverageCandidates  float64       `json:"average_candidates"`
	RemainingSolutions string        `json:"remaining_solutions"`
	UniqueSolution     bool          `json:"unique_solution"`
	NextCell           *CellAnalysis `json:"next_cell,omitempty"`
}

// SolveStep records one probability-guided choice.
type SolveStep struct {
	Step                     int          `json:"step"`
	RemainingSolutionsBefore string       `json:"remaining_solutions_before"`
	RemainingSolutionsAfter  string       `json:"remaining_solutions_after"`
	Cell                     CellAnalysis `json:"cell"`
}

// SolveResult records a complete probability-guided solve.
type SolveResult struct {
	InitialSolutions string      `json:"initial_solutions"`
	Steps            []SolveStep `json:"steps"`
	FinalGrid        string      `json:"final_grid"`
	Solved           bool        `json:"solved"`
}

// Analyze counts every valid completion compatible with the current grid and
// computes the exact candidate distribution for the next empty cell (row-major).
//
// Candidate probability is:
//
//	completions after choosing the candidate / all remaining completions
//
// Counts use arbitrary-precision integers so the result does not overflow.
func Analyze(grid Grid) (Analysis, error) {
	counter := newExactCounter()
	return analyzeWithCounter(grid, counter)
}

// CountSolutions returns the exact number of valid complete Sudokus compatible
// with the current grid.
func CountSolutions(grid Grid) (*big.Int, error) {
	if err := Validate(grid); err != nil {
		return nil, err
	}
	counter := newExactCounter()
	return counter.count(grid), nil
}

// SolveByProbability repeatedly analyzes the next empty cell and selects the
// digit contained in the largest share of the remaining exact solution space.
// On a uniquely solvable Sudoku every selected digit is guaranteed (100%).
func SolveByProbability(grid Grid) (SolveResult, error) {
	if err := Validate(grid); err != nil {
		return SolveResult{}, err
	}

	counter := newExactCounter()
	initial, err := analyzeWithCounter(grid, counter)
	if err != nil {
		return SolveResult{}, err
	}
	result := SolveResult{InitialSolutions: initial.RemainingSolutions}
	if initial.RemainingSolutions == "0" {
		return result, nil
	}

	current := grid
	for step := 1; ; step++ {
		analysis, err := analyzeWithCounter(current, counter)
		if err != nil {
			return SolveResult{}, err
		}
		if analysis.NextCell == nil {
			result.FinalGrid = Compact(current)
			result.Solved = Validate(current) == nil && analysis.EmptyCells == 0 && analysis.RemainingSolutions == "1"
			return result, nil
		}

		cell := *analysis.NextCell
		idx := (cell.Row-1)*Size + (cell.Column - 1)
		if cell.RecommendedDigit == 0 || cell.RecommendedRemainingSolutions == "0" {
			return SolveResult{}, fmt.Errorf("analysis produced no viable candidate for r%dc%d", cell.Row, cell.Column)
		}

		result.Steps = append(result.Steps, SolveStep{
			Step:                     step,
			RemainingSolutionsBefore: analysis.RemainingSolutions,
			RemainingSolutionsAfter:  cell.RecommendedRemainingSolutions,
			Cell:                     cell,
		})
		current[idx] = cell.RecommendedDigit
	}
}

func analyzeWithCounter(grid Grid, counter *exactCounter) (Analysis, error) {
	if err := Validate(grid); err != nil {
		return Analysis{}, err
	}

	analysis, _ := summarize(grid)
	total := counter.count(grid)
	analysis.RemainingSolutions = total.String()
	analysis.UniqueSolution = total.Cmp(big.NewInt(1)) == 0

	if total.Sign() == 0 || analysis.EmptyCells == 0 {
		return analysis, nil
	}

	idx := firstEmptyCell(grid)
	cell := analyzeCell(grid, idx, total, counter)
	analysis.NextCell = &cell
	return analysis, nil
}

func analyzeCell(grid Grid, idx int, total *big.Int, counter *exactCounter) CellAnalysis {
	cell := CellAnalysis{
		Row:        idx/Size + 1,
		Column:     idx%Size + 1,
		Candidates: make([]CandidateAnalysis, 0, Size),
	}
	allowed := candidateMask(grid, idx)
	best := new(big.Int)

	for digit := uint8(1); digit <= 9; digit++ {
		count := new(big.Int)
		isAllowed := allowed&(1<<digit) != 0
		if isAllowed {
			next := grid
			next[idx] = digit
			count = counter.count(next)
		}
		candidate := CandidateAnalysis{
			Digit:              digit,
			LocallyAllowed:     isAllowed,
			RemainingSolutions: count.String(),
			ProbabilityPercent: solutionPercent(count, total),
		}
		cell.Candidates = append(cell.Candidates, candidate)

		if count.Cmp(best) > 0 {
			best.Set(count)
			cell.RecommendedDigit = digit
			cell.RecommendedRemainingSolutions = count.String()
			cell.RecommendedProbabilityPercent = candidate.ProbabilityPercent
		}
	}

	cell.Guaranteed = best.Cmp(total) == 0 && total.Sign() > 0
	return cell
}

func solutionPercent(part, total *big.Int) float64 {
	if total.Sign() == 0 || part.Sign() == 0 {
		return 0
	}
	ratio := new(big.Rat).SetFrac(part, total)
	value, _ := ratio.Float64()
	return value * 100
}

func firstEmptyCell(grid Grid) int {
	for idx, value := range grid {
		if value == 0 {
			return idx
		}
	}
	return -1
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
			result.MinCandidates = 0
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

type exactCounter struct {
	memo map[Grid]*big.Int
}

func newExactCounter() *exactCounter {
	return &exactCounter{memo: make(map[Grid]*big.Int)}
}

func (c *exactCounter) count(start Grid) *big.Int {
	if cached, ok := c.memo[start]; ok {
		return new(big.Int).Set(cached)
	}

	grid := start
	if !propagateSingles(&grid) {
		zero := new(big.Int)
		c.memo[start] = zero
		return new(big.Int)
	}
	if cached, ok := c.memo[grid]; ok {
		c.memo[start] = new(big.Int).Set(cached)
		return new(big.Int).Set(cached)
	}

	bestIdx := -1
	bestMask := uint16(0)
	bestCount := 10
	for idx, value := range grid {
		if value != 0 {
			continue
		}
		mask := candidateMask(grid, idx)
		count := bits.OnesCount16(mask)
		if count < bestCount {
			bestIdx = idx
			bestMask = mask
			bestCount = count
		}
	}

	if bestIdx == -1 {
		one := big.NewInt(1)
		c.memo[grid] = new(big.Int).Set(one)
		c.memo[start] = new(big.Int).Set(one)
		return one
	}

	total := new(big.Int)
	for _, digit := range digitsFromMask(bestMask) {
		next := grid
		next[bestIdx] = digit
		total.Add(total, c.count(next))
	}

	c.memo[grid] = new(big.Int).Set(total)
	c.memo[start] = new(big.Int).Set(total)
	return new(big.Int).Set(total)
}

// HasSolution checks whether at least one valid completion exists without
// counting every completion.
func HasSolution(grid Grid) bool {
	if err := Validate(grid); err != nil || !unitsViable(grid) {
		return false
	}
	copyGrid := grid
	return searchAny(&copyGrid)
}

func searchAny(grid *Grid) bool {
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
		if searchAny(&next) {
			return true
		}
	}
	return false
}
