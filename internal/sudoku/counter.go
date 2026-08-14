package sudoku

import (
	"math/big"
	"math/bits"
)

const (
	exactMemoSmallSlots        = 128
	exactMemoLargeSlots        = 256
	probabilityMemoSlots       = 256
	exactMemoPressureCutoff    = 190
	memoMaximumLoadNumerator   = 3
	memoMaximumLoadDenominator = 4
)

// exactCount keeps the overwhelmingly common 0/1/small counts inline and only
// allocates a big.Int when the exact result genuinely exceeds uint64.
type exactCount struct {
	small uint64
	large *big.Int
}

func zeroExactCount() exactCount {
	return exactCount{}
}

func oneExactCount() exactCount {
	return exactCount{small: 1}
}

func (count exactCount) bigInt() *big.Int {
	if count.large != nil {
		return new(big.Int).Set(count.large)
	}
	return new(big.Int).SetUint64(count.small)
}

func (count *exactCount) add(other exactCount) {
	if other.large == nil && other.small == 0 {
		return
	}
	if count.large != nil {
		if other.large != nil {
			count.large.Add(count.large, other.large)
			return
		}
		var value big.Int
		value.SetUint64(other.small)
		count.large.Add(count.large, &value)
		return
	}
	if other.large != nil {
		count.large = new(big.Int).SetUint64(count.small)
		count.large.Add(count.large, other.large)
		count.small = 0
		return
	}

	sum, carry := bits.Add64(count.small, other.small, 0)
	if carry == 0 {
		count.small = sum
		return
	}

	promoted := new(big.Int).SetUint64(count.small)
	var otherValue big.Int
	otherValue.SetUint64(other.small)
	promoted.Add(promoted, &otherValue)
	count.small = 0
	count.large = promoted
}

// packedGrid is an exact 4-bit-per-cell key for memoization. Six uint64 words
// hold all 81 cells (324 bits) without keeping an 81-byte Grid in every slot.
type packedGrid [6]uint64

func packGrid(grid Grid) packedGrid {
	var packed packedGrid
	for idx, digit := range grid {
		if digit != 0 {
			packed.set(idx, digit)
		}
	}
	return packed
}

func (packed *packedGrid) set(idx int, digit uint8) {
	word := idx / 16
	shift := uint((idx % 16) * 4)
	mask := uint64(0xF) << shift
	packed[word] = (packed[word] &^ mask) | (uint64(digit) << shift)
}

type memoEntry struct {
	key   packedGrid
	value exactCount
	hash  uint64
	used  bool
}

type memoTable struct {
	entries []memoEntry
	used    int
	limit   int
}

func newMemoTable(slots int) memoTable {
	if slots < 1 {
		return memoTable{}
	}
	powerOfTwo := 1
	for powerOfTwo < slots {
		powerOfTwo <<= 1
	}
	return memoTable{
		entries: make([]memoEntry, powerOfTwo),
		limit:   powerOfTwo * memoMaximumLoadNumerator / memoMaximumLoadDenominator,
	}
}

func hashPackedGrid(key packedGrid) uint64 {
	hash := uint64(0x9e3779b97f4a7c15)
	for _, word := range key {
		word ^= word >> 30
		word *= 0xbf58476d1ce4e5b9
		word ^= word >> 27
		word *= 0x94d049bb133111eb
		word ^= word >> 31
		hash ^= word
		hash *= 0x9e3779b97f4a7c15
	}
	return hash
}

func (memo *memoTable) get(key packedGrid) (exactCount, bool) {
	if len(memo.entries) == 0 {
		return exactCount{}, false
	}
	hash := hashPackedGrid(key)
	mask := uint64(len(memo.entries) - 1)
	idx := hash & mask
	for probe := 0; probe < len(memo.entries); probe++ {
		entry := &memo.entries[idx]
		if !entry.used {
			return exactCount{}, false
		}
		if entry.hash == hash && entry.key == key {
			return entry.value, true
		}
		idx = (idx + 1) & mask
	}
	return exactCount{}, false
}

func (memo *memoTable) set(key packedGrid, value exactCount) {
	if len(memo.entries) == 0 {
		return
	}
	hash := hashPackedGrid(key)
	mask := uint64(len(memo.entries) - 1)
	idx := hash & mask
	for probe := 0; probe < len(memo.entries); probe++ {
		entry := &memo.entries[idx]
		if !entry.used {
			if memo.used >= memo.limit {
				return
			}
			entry.used = true
			entry.hash = hash
			entry.key = key
			entry.value = value
			memo.used++
			return
		}
		if entry.hash == hash && entry.key == key {
			entry.value = value
			return
		}
		idx = (idx + 1) & mask
	}
}

type searchFrame struct {
	state         searchState
	start         packedGrid
	remainingMask uint16
	branchIdx     uint8
	total         exactCount
	initialized   bool
}

type exactCounter struct {
	memo           memoTable
	stack          []searchFrame
	degreeTieBreak bool
}

func newExactCounter() *exactCounter {
	return &exactCounter{degreeTieBreak: true}
}

func newProbabilityCounter() *exactCounter {
	return &exactCounter{}
}

func (counter *exactCounter) count(start Grid) *big.Int {
	return counter.countExact(start).bigInt()
}

func (counter *exactCounter) countExact(start Grid) exactCount {
	root := newSearchState(start)
	counter.prepareWorkspace(&root)
	counter.stack[0] = searchFrame{state: root}
	depth := 0

	for {
		frame := &counter.stack[depth]
		var result exactCount
		finished := false

		if !frame.initialized {
			frame.start = frame.state.key
			if cached, ok := counter.memo.get(frame.start); ok {
				result = cached
				finished = true
			} else if !propagateStateSingles(&frame.state) {
				result = zeroExactCount()
				counter.memo.set(frame.start, result)
				finished = true
			} else if cached, ok := counter.memo.get(frame.state.key); ok {
				result = cached
				counter.memo.set(frame.start, cached)
				finished = true
			} else {
				bestIdx, bestMask := selectBranchCellRowMajor(&frame.state)
				if counter.degreeTieBreak {
					bestIdx, bestMask = selectBranchCell(&frame.state)
				}
				if bestIdx == -1 {
					result = oneExactCount()
					counter.memo.set(frame.state.key, result)
					counter.memo.set(frame.start, result)
					finished = true
				} else {
					frame.branchIdx = uint8(bestIdx)
					frame.remainingMask = bestMask
					frame.total = zeroExactCount()
					frame.initialized = true
				}
			}
		}

		if !finished {
			if frame.remainingMask != 0 {
				bit := frame.remainingMask & (^frame.remainingMask + 1)
				digit := uint8(bits.TrailingZeros16(bit))
				frame.remainingMask &^= bit
				next := frame.state
				next.place(int(frame.branchIdx), digit)
				depth++
				counter.stack[depth] = searchFrame{state: next}
				continue
			}

			result = frame.total
			counter.memo.set(frame.state.key, result)
			counter.memo.set(frame.start, result)
		}

		if depth == 0 {
			return result
		}
		counter.stack[depth] = searchFrame{}
		depth--
		counter.stack[depth].total.add(result)
	}
}

func (counter *exactCounter) prepareWorkspace(root *searchState) {
	requiredFrames := 1
	for _, digit := range root.grid {
		if digit == 0 {
			requiredFrames++
		}
	}
	if len(counter.stack) < requiredFrames {
		counter.stack = make([]searchFrame, requiredFrames)
	}
	if len(counter.memo.entries) == 0 {
		counter.memo = newMemoTable(counter.memoWindowSlots(root))
	}
}

// memoWindowSlots intentionally stays small. The exact number of reachable
// states cannot be known without performing the search first. A bounded cache
// avoids map growth and GC pressure; when full, new states are simply computed
// without caching, which changes performance but never the exact answer.
func (counter *exactCounter) memoWindowSlots(root *searchState) int {
	if !counter.degreeTieBreak {
		return probabilityMemoSlots
	}
	pressure := 0
	for idx, digit := range root.grid {
		if digit != 0 {
			continue
		}
		candidateCount := bits.OnesCount16(root.candidateMask(idx))
		if candidateCount > 1 {
			pressure += candidateCount - 1
		}
	}
	if pressure < exactMemoPressureCutoff {
		return exactMemoSmallSlots
	}
	return exactMemoLargeSlots
}
