package pack_test

import (
	"context"
	"testing"

	"pfg/internal/pack"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	sizes []int
}

func (m *mockRepo) GetPackSizes(ctx context.Context) ([]int, error) {
	return m.sizes, nil
}

func (m *mockRepo) InsertPackSize(ctx context.Context, size int) error {
	m.sizes = append(m.sizes, size)
	return nil
}

func (m *mockRepo) DeletePackSize(ctx context.Context, size int) error {
	result := make([]int, 0)
	for _, s := range m.sizes {
		if s != size {
			result = append(result, s)
		}
	}
	m.sizes = result
	return nil
}

func TestCalculate(t *testing.T) {
	repo := &mockRepo{sizes: []int{250, 500, 1000, 2000, 5000}}
	service := pack.NewService(repo)

	tests := []struct {
		name     string
		quantity int
		expected pack.PackResult
	}{
		{
			name:     "Exact Pack",
			quantity: 250,
			expected: pack.PackResult{TotalItems: 250, TotalPacks: 1, Packs: map[int]int{250: 1}},
		},
		{
			name:     "Minimal Overage Fewest Packs",
			quantity: 251,
			expected: pack.PackResult{TotalItems: 500, TotalPacks: 1, Packs: map[int]int{500: 1}},
		},
		{
			name:     "Avoid Extra Items",
			quantity: 501,
			expected: pack.PackResult{TotalItems: 750, TotalPacks: 2, Packs: map[int]int{500: 1, 250: 1}},
		},
		{
			name:     "Multiple Large Packs",
			quantity: 12001,
			expected: pack.PackResult{TotalItems: 12250, TotalPacks: 4, Packs: map[int]int{5000: 2, 2000: 1, 250: 1}},
		},
		{
			name:     "Small Order",
			quantity: 1,
			expected: pack.PackResult{TotalItems: 250, TotalPacks: 1, Packs: map[int]int{250: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Calculate(context.Background(), tt.quantity)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.TotalItems, result.TotalItems)
			assert.Equal(t, tt.expected.TotalPacks, result.TotalPacks)
			assert.Equal(t, tt.expected.Packs, result.Packs)
		})
	}
}
