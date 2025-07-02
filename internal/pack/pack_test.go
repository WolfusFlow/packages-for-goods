package pack_test

import (
	"pfg/internal/pack"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactPack(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	result := pack.CalculateOptimalPacks(250, sizes)

	assert.Equal(t, 250, result.TotalItems)
	assert.Equal(t, 1, result.TotalPacks)
	assert.Equal(t, map[int]int{250: 1}, result.Packs)
}

func TestMinimalOverageAndFewerPacks(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	result := pack.CalculateOptimalPacks(251, sizes)

	assert.Equal(t, 500, result.TotalItems)
	assert.Equal(t, 1, result.TotalPacks)
	assert.Equal(t, map[int]int{500: 1}, result.Packs)
}

func TestAvoidExtraItems(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	result := pack.CalculateOptimalPacks(501, sizes)

	assert.Equal(t, 750, result.TotalItems)
	assert.Equal(t, 2, result.TotalPacks)
	assert.Equal(t, map[int]int{500: 1, 250: 1}, result.Packs)
}

func TestMultipleLargePacks(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	result := pack.CalculateOptimalPacks(12001, sizes)

	assert.Equal(t, 12250, result.TotalItems)
	assert.Equal(t, 4, result.TotalPacks)
	assert.Equal(t, map[int]int{5000: 2, 2000: 1, 250: 1}, result.Packs)
}

func TestSmallOrder(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	result := pack.CalculateOptimalPacks(1, sizes)

	assert.Equal(t, 250, result.TotalItems)
	assert.Equal(t, 1, result.TotalPacks)
	assert.Equal(t, map[int]int{250: 1}, result.Packs)
}
