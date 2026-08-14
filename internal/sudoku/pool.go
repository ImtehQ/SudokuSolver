package sudoku

import "sync"

var probabilityCounterPool = sync.Pool{
	New: func() any {
		return newProbabilityCounter()
	},
}

func acquireProbabilityCounter() *exactCounter {
	return probabilityCounterPool.Get().(*exactCounter)
}

func releaseProbabilityCounter(counter *exactCounter) {
	counter.resetWorkspace()
	probabilityCounterPool.Put(counter)
}

// resetWorkspace keeps the allocated memo and DFS windows but removes every
// puzzle-specific key, count, and search frame before the counter is reused.
func (counter *exactCounter) resetWorkspace() {
	for idx := range counter.memo.entries {
		counter.memo.entries[idx] = memoEntry{}
	}
	counter.memo.used = 0
	for idx := range counter.stack {
		counter.stack[idx] = searchFrame{}
	}
}
