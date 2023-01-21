package matrix

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/tuple"
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

func TestMultiplyByMatrix(t *testing.T) {
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

func TestMultipleByTuple(t *testing.T) {
	m := NewFromSlice([][]float64{
		{1, 2, 3, 4},
		{2, 4, 4, 2},
		{8, 6, 4, 1},
		{0, 0, 0, 1},
	})

	tu := tuple.New(1, 2, 3, 1)

	product := m.MulTuple(tu)

	assert.Equal(t, tuple.New(18, 24, 33, 1), product)
}

func TestMultiplyMatrixByIdentityMatrix(t *testing.T) {
	m := NewFromSlice([][]float64{
		{0, 1, 2, 4},
		{1, 2, 4, 8},
		{2, 4, 8, 16},
		{4, 8, 16, 32},
	})

	assert.Equal(t, m, m.Mul(Identity4()))
}

func TestMultiplyIdentityMatrixByTuple(t *testing.T) {
	tu := tuple.New(1, 2, 3, 4)

	assert.Equal(t, tu, Identity4().MulTuple(tu))
}

func TestTranspose(t *testing.T) {
	m := NewFromSlice([][]float64{
		{0, 9, 3, 0},
		{9, 8, 0, 8},
		{1, 8, 5, 3},
		{0, 0, 5, 8},
	})

	expected := NewFromSlice([][]float64{
		{0, 9, 1, 0},
		{9, 8, 8, 0},
		{3, 0, 5, 5},
		{0, 8, 3, 8},
	})

	assert.Equal(t, expected, m.Transpose())
}

func TestTransposeIdentity(t *testing.T) {
	assert.Equal(t, Identity4(), Identity4().Transpose())
}

func Test2x2Determinant(t *testing.T) {
	m := NewFromSlice([][]float64{
		{1, 5},
		{-3, 2},
	})

	assert.Equal(t, 17.0, m.Determinant())
}

func Test3x3Submatrix(t *testing.T) {
	m := NewFromSlice([][]float64{
		{1, 5, 0},
		{-3, 2, 7},
		{0, 6, -3},
	})

	expected := NewFromSlice([][]float64{
		{-3, 2},
		{0, 6},
	})

	assert.Equal(t, expected, m.Submatrix(0, 2))
}

func Test4x4Submatrix(t *testing.T) {
	m := NewFromSlice([][]float64{
		{-6, 1, 1, 6},
		{-8, 5, 8, 6},
		{-1, 0, 8, 2},
		{-7, 1, -1, 1},
	})

	expected := NewFromSlice([][]float64{
		{-6, 1, 6},
		{-8, 8, 6},
		{-7, -1, 1},
	})

	assert.Equal(t, expected, m.Submatrix(2, 1))
}

func Test3x3Minor(t *testing.T) {
	a := NewFromSlice([][]float64{
		{3, 5, 0},
		{2, -1, -7},
		{6, -1, 5},
	})

	b := a.Submatrix(1, 0)

	assert.Equal(t, 25.0, b.Determinant())
	assert.Equal(t, 25.0, a.Minor(1, 0))
}

func Test3x3Cofactor(t *testing.T) {
	a := NewFromSlice([][]float64{
		{3, 5, 0},
		{2, -1, -7},
		{6, -1, 5},
	})

	assert.Equal(t, -12.0, a.Minor(0, 0))
	assert.Equal(t, -12.0, a.Cofactor(0, 0))
	assert.Equal(t, 25.0, a.Minor(1, 0))
	assert.Equal(t, -25.0, a.Cofactor(1, 0))
}

func Test3x3Determinant(t *testing.T) {
	a := NewFromSlice([][]float64{
		{1, 2, 6},
		{-5, 8, -4},
		{2, 6, 4},
	})

	assert.Equal(t, 56.0, a.Cofactor(0, 0))
	assert.Equal(t, 12.0, a.Cofactor(0, 1))
	assert.Equal(t, -46.0, a.Cofactor(0, 2))
	assert.Equal(t, -196.0, a.Determinant())
}

func Test4x4Determinant(t *testing.T) {
	a := NewFromSlice([][]float64{
		{-2, -8, 3, 5},
		{-3, 1, 7, 3},
		{1, 2, -9, 6},
		{-6, 7, 7, -9},
	})

	assert.Equal(t, 690.0, a.Cofactor(0, 0))
	assert.Equal(t, 447.0, a.Cofactor(0, 1))
	assert.Equal(t, 210.0, a.Cofactor(0, 2))
	assert.Equal(t, 51.0, a.Cofactor(0, 3))
	assert.Equal(t, -4071.0, a.Determinant())
}

func TestInvertible(t *testing.T) {
	a := NewFromSlice([][]float64{
		{6, 4, 4, 4},
		{5, 5, 7, 6},
		{4, -9, 3, -7},
		{9, 1, 7, -6},
	})

	assert.Equal(t, -2120.0, a.Determinant())
	assert.True(t, a.IsInvertible())
}

func TestNotInvertible(t *testing.T) {
	a := NewFromSlice([][]float64{
		{-4, 2, -2, -3},
		{9, 6, 2, 6},
		{0, -5, 1, -5},
		{0, 0, 0, 0},
	})

	assert.Equal(t, 0.0, a.Determinant())
	assert.False(t, a.IsInvertible())
}

func TestInverse(t *testing.T) {
	a := NewFromSlice([][]float64{
		{-5, 2, 6, -8},
		{1, -5, 1, 8},
		{7, 7, -6, -7},
		{1, -3, 7, 4},
	})

	b := a.Inverse()

	assert.Equal(t, 532.0, a.Determinant())
	assert.Equal(t, -160.0, a.Cofactor(2, 3))
	assert.Equal(t, -160.0/532.0, b.At(3, 2))
	assert.Equal(t, 105.0, a.Cofactor(3, 2))
	assert.Equal(t, 105.0/532.0, b.At(2, 3))

	expectedInverse := NewFromSlice([][]float64{
		{0.21805, 0.45113, 0.24060, -0.04511},
		{-0.80827, -1.45677, -0.44361, 0.52068},
		{-0.07895, -0.22368, -0.05263, 0.19737},
		{-0.52256, -0.81391, -0.30075, 0.30639},
	})

	assert.True(t, expectedInverse.Equals(b))
}

func TestInverse2(t *testing.T) {
	a := NewFromSlice([][]float64{
		{8, -5, 9, 2},
		{7, 5, 6, 1},
		{-6, 0, 9, 6},
		{-3, 0, -9, -4},
	})

	b := a.Inverse()

	expectedInverse := NewFromSlice([][]float64{
		{-0.15385, -0.15385, -0.28205, -0.53846},
		{-0.07692, 0.12308, 0.02564, 0.03077},
		{0.35897, 0.35897, 0.43590, 0.92308},
		{-0.69231, -0.69231, -0.76923, -1.92308},
	})

	assert.True(t, expectedInverse.Equals(b))
}

func TestInverse3(t *testing.T) {
	a := NewFromSlice([][]float64{
		{9, 3, 0, 9},
		{-5, -2, -6, -3},
		{-4, 9, 6, 4},
		{-7, 6, 6, 2},
	})

	b := a.Inverse()

	expectedInverse := NewFromSlice([][]float64{
		{-0.04074, -0.07778, 0.14444, -0.22222},
		{-0.07778, 0.03333, 0.36667, -0.33333},
		{-0.02901, -0.14630, -0.10926, 0.12963},
		{0.17778, 0.06667, -0.26667, 0.33333},
	})

	assert.True(t, expectedInverse.Equals(b))
}

func TestMultiplyProductByItsInverse(t *testing.T) {
	a := NewFromSlice([][]float64{
		{3, -9, 7, 3},
		{3, -8, 2, -9},
		{-4, 4, 4, 1},
		{-6, 5, -1, 1},
	})

	b := NewFromSlice([][]float64{
		{8, 2, 2, 2},
		{3, -1, 7, 0},
		{7, 0, 5, 4},
		{6, -2, 0, 5},
	})

	assert.True(t, a.Equals(a.Mul(b).Mul(b.Inverse())))
}

func TestInverseOfIdentityMatrix(t *testing.T) {
	assert.Equal(t, Identity4(), Identity4().Inverse())
}

func TestMultiplyMatrixByItsInverse(t *testing.T) {
	a := NewFromSlice([][]float64{
		{3, -9, 7, 3},
		{3, -8, 2, -9},
		{-4, 4, 4, 1},
		{-6, 5, -1, 1},
	})

	assert.True(t, Identity4().Equals(a.Mul(a.Inverse())))
}

func TestInverseTransposeCommutativity(t *testing.T) {
	a := NewFromSlice([][]float64{
		{3, -9, 7, 3},
		{3, -8, 2, -9},
		{-4, 4, 4, 1},
		{-6, 5, -1, 1},
	})

	assert.True(t, a.Inverse().Transpose().Equals(a.Transpose().Inverse()))
}

func TestCopy(t *testing.T) {
	a := NewFromSlice([][]float64{
		{3, -9, 7, 3},
		{3, -8, 2, -9},
		{-4, 4, 4, 1},
		{-6, 5, -1, 1},
	})

	assert.Equal(t, a, a.Copy())
}
