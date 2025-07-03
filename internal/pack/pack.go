package pack

import (
	"context"
	"errors"
	"sort"
)

type PackResult struct {
	TotalItems int
	TotalPacks int
	Packs      map[int]int
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListPacks(ctx context.Context) ([]int, error) {
	return s.repo.GetPackSizes(ctx)
}

func (s *Service) AddPack(ctx context.Context, size int) error {
	if size <= 0 {
		return errors.New("invalid pack size")
	}
	return s.repo.InsertPackSize(ctx, size)
}

func (s *Service) RemovePack(ctx context.Context, size int) error {
	return s.repo.DeletePackSize(ctx, size)
}

func (s *Service) Calculate(ctx context.Context, quantity int) (PackResult, error) {
	sizes, err := s.repo.GetPackSizes(ctx)
	if err != nil {
		return PackResult{}, err
	}

	if len(sizes) == 0 {
		return PackResult{}, errors.New("no pack sizes available")
	}

	sort.Ints(sizes)

	maxSize := sizes[len(sizes)-1]
	limit := quantity + maxSize*2

	type state struct {
		packCount   int
		combination map[int]int
	}

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

	for i := quantity; i <= limit; i++ {
		if dp[i] != nil {
			return PackResult{
				TotalItems: i,
				TotalPacks: dp[i].packCount,
				Packs:      dp[i].combination,
			}, nil
		}
	}

	return PackResult{}, errors.New("no valid pack combination found")
}

func copyMap(m map[int]int) map[int]int {
	cp := make(map[int]int, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
