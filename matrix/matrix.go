package matrix

import (
	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Matrix struct {
	m [][]float64
}

func New(w, h int) Matrix {
	m := make([][]float64, h)
	for i := range m {
		m[i] = make([]float64, w)
	}
	return Matrix{m}
}

func NewFromSlice(s [][]float64) Matrix {
	if len(s) == 0 {
		panic("empty matrix")
	}
	width := len(s[0])
	for i := range s {
		if len(s[i]) != width {
			panic("invalid matrix dimensions")
		}
	}
	return Matrix{s}
}

func (m Matrix) Copy() Matrix {
	rows := len(m.m)
	cols := len(m.m[0])
	m2 := make([][]float64, rows)
	data := make([]float64, rows*cols)
	for i := range m.m {
		start := i * cols
		end := start + cols
		m2[i] = data[start:end:end]
		copy(m2[i], m.m[i])
	}
	return NewFromSlice(m2)
}

func (m Matrix) At(x, y int) float64 {
	return m.m[x][y]
}

func (m Matrix) Equals(m2 Matrix) bool {
	if len(m.m) != len(m2.m) {
		return false
	}
	if len(m.m) == 0 {
		return true
	}
	if len(m.m[0]) != len(m2.m[0]) {
		return false
	}

	for y := range m.m {
		for x := range m.m[y] {
			if !float.Equal(m.m[x][y], m2.m[x][y]) {
				return false
			}
		}
	}

	return true
}

func (m Matrix) Mul(m2 Matrix) Matrix {
	if len(m.m[0]) != len(m2.m) {
		panic("invalid dimensions for matrix multiplication")
	}
	product := New(len(m2.m[0]), len(m.m))

	for y := range product.m {
		for x := range product.m[y] {
			sum := 0.0
			for k := 0; k < len(m2.m); k++ {
				sum += m.m[x][k] * m2.m[k][y]
			}
			product.m[x][y] = sum
		}
	}

	return product
}

func (m Matrix) MulTuple(t tuple.Tuple) tuple.Tuple {
	if len(m.m) != 4 || len(m.m[0]) != 4 {
		panic("invalid dimensions for matrix-tuple multiplication")
	}

	x := m.m[0][0]*t.X + m.m[0][1]*t.Y + m.m[0][2]*t.Z + m.m[0][3]*t.W
	y := m.m[1][0]*t.X + m.m[1][1]*t.Y + m.m[1][2]*t.Z + m.m[1][3]*t.W
	z := m.m[2][0]*t.X + m.m[2][1]*t.Y + m.m[2][2]*t.Z + m.m[2][3]*t.W
	w := m.m[3][0]*t.X + m.m[3][1]*t.Y + m.m[3][2]*t.Z + m.m[3][3]*t.W

	return tuple.New(x, y, z, w)
}

func (m Matrix) Transpose() Matrix {
	t := New(len(m.m), len(m.m[0]))

	for j := 0; j < len(m.m); j++ {
		for i := 0; i < len(m.m[0]); i++ {
			t.m[i][j] = m.m[j][i]
		}
	}

	return t
}

func (m Matrix) Determinant() float64 {
	if len(m.m) == 2 {
		return m.m[0][0]*m.m[1][1] - m.m[0][1]*m.m[1][0]
	}

	sum := 0.0

	for j := range m.m {
		cofactor := m.Cofactor(0, j)
		sum += m.m[0][j] * cofactor
	}

	return sum
}

func (m Matrix) Submatrix(row, col int) Matrix {
	sub := New(len(m.m)-1, len(m.m[0])-1)

	for i := 0; i < len(m.m); i++ {
		for j := 0; j < len(m.m[0]); j++ {
			if i != row && j != col {
				newi := i
				newj := j
				if i > row {
					newi -= 1
				}
				if j > col {
					newj -= 1
				}

				sub.m[newi][newj] = m.m[i][j]
			}
		}
	}

	return sub
}

func (m Matrix) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

func (m Matrix) Cofactor(row, col int) float64 {
	minor := m.Minor(row, col)
	if (row+col)%2 != 0 {
		return -minor
	}
	return minor
}

func (m Matrix) IsInvertible() bool {
	return m.Determinant() != 0
}

func (m Matrix) Inverse() Matrix {
	determinant := m.Determinant()
	if determinant == 0 {
		panic("matrix not invertible")
	}

	inverse := New(len(m.m), len(m.m[0]))
	for i := 0; i < len(m.m); i++ {
		for j := 0; j < len(m.m[0]); j++ {
			c := m.Cofactor(i, j)
			inverse.m[j][i] = c / determinant
		}
	}

	return inverse
}

func Identity4() Matrix {
	return NewFromSlice([][]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	})
}
