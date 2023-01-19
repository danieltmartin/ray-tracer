package matrix

import "github.com/danieltmartin/ray-tracer/float"

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
