package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMatrix(t *testing.T) {
	m := New(4, 4)

	assert.Equal(t, 0.0, m.At(0, 0))
	assert.Equal(t, 0.0, m.At(3, 3))
}

func TestCreate4x4MatrixFromSlize(t *testing.T) {
	nums := [][]float64{
		{1, 2, 3, 4},
		{5.5, 6.5, 7.5, 8.5},
		{9, 10, 11, 12},
		{13.5, 14.5, 15.5, 16.5},
	}

	m := NewFromSlice(nums)

	assert.Equal(t, 1.0, m.At(0, 0))
	assert.Equal(t, 7.5, m.At(1, 2))
}

func TestCreate2x2MatrixFromSlice(t *testing.T) {
	nums := [][]float64{
		{-3, 5},
		{1, -2},
	}

	m := NewFromSlice(nums)

	assert.Equal(t, -3.0, m.At(0, 0))
	assert.Equal(t, 5.0, m.At(0, 1))
	assert.Equal(t, 1.0, m.At(1, 0))
	assert.Equal(t, -2.0, m.At(1, 1))
}

func TestExactlyEqual(t *testing.T) {
	m1 := NewFromSlice([][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 8, 7, 6},
		{5, 4, 3, 2},
	})

	assert.True(t, m1.Equals(m1))
}

func TestNotEquals(t *testing.T) {
	m1 := NewFromSlice([][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 8, 7, 6},
		{5, 4, 3, 2},
	})

	m2 := NewFromSlice([][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 8, 7, 0},
		{5, 4, 3, 2},
	})

	assert.False(t, m1.Equals(m2))
}

func TestApproxEquals(t *testing.T) {
	m1 := NewFromSlice([][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 8, 7, 6},
		{5, 4, 3, 2},
	})

	m2 := NewFromSlice([][]float64{
		{1.000001, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 8, 7, 6},
		{5, 4, 3, 2},
	})

	assert.True(t, m1.Equals(m2))
}

func TestMultiply(t *testing.T) {
	m1 := NewFromSlice([][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 8, 7, 6},
		{5, 4, 3, 2},
	})

	m2 := NewFromSlice([][]float64{
		{-2, 1, 2, 3},
		{3, 2, 1, -1},
		{4, 3, 6, 5},
		{1, 2, 7, 8},
	})

	product := NewFromSlice([][]float64{
		{20, 22, 50, 48},
		{44, 54, 114, 108},
		{40, 58, 110, 102},
		{16, 26, 46, 42},
	})

	assert.Equal(t, product, m1.Mul(m2))
}
