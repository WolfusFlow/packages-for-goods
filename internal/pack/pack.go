package pack

import (
	"sort"
)

type PackResult struct {
	TotalItems int
	TotalPacks int
	Packs      map[int]int
}

// state struct for pack count and combination for total i
type state struct {
	packCount   int
	combination map[int]int
}

// CalculateOptimalPacks finds the optimal combination.
func CalculateOptimalPacks(quantity int, sizes []int) PackResult {
	sort.Ints(sizes)

	maxSize := sizes[len(sizes)-1]
	limit := quantity + maxSize*2 // over-allocate space for overage options

	dp := make([]*state, limit+1)
	dp[0] = &state{packCount: 0, combination: map[int]int{}}

	for i := 0; i <= limit; i++ {
		if dp[i] == nil {
			continue
		}
		for _, size := range sizes {
			next := i + size
			if next > limit {
				continue
			}
			newCount := dp[i].packCount + 1
			if dp[next] == nil || newCount < dp[next].packCount {
				newComb := copyMap(dp[i].combination)
				newComb[size]++
				dp[next] = &state{packCount: newCount, combination: newComb}
			}
		}
	}

	// Find the best valid result starting from quantity upward
	for i := quantity; i <= limit; i++ {
		if dp[i] != nil {
			return PackResult{
				TotalItems: i,
				TotalPacks: dp[i].packCount,
				Packs:      dp[i].combination,
			}
		}
	}

	// Fallback. Should never happen
	return PackResult{
		TotalItems: 0,
		TotalPacks: 0,
		Packs:      map[int]int{},
	}
}

func copyMap(m map[int]int) map[int]int {
	cp := make(map[int]int, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
