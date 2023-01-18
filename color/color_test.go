package color

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/stretchr/testify/assert"
)

func TestNewColor(t *testing.T) {
	c := New(-0.5, 0.4, 1.7)

	assert.Equal(t, c.R(), -0.5)
	assert.Equal(t, c.G(), 0.4)
	assert.Equal(t, c.B(), 1.7)
}

func TestAddColors(t *testing.T) {
	c1 := New(0.9, 0.6, 0.75)
	c2 := New(0.7, 0.1, 0.25)

	assert.Equal(t, c1.Add(c2), New(1.6, 0.7, 1.0))
}

func TestSubtractColors(t *testing.T) {
	c1 := New(0.9, 0.6, 0.75)
	c2 := New(0.7, 0.1, 0.25)

	assertAlmost(t, c1.Sub(c2), New(0.2, 0.5, 0.5))
}

func TestMultiplyColors(t *testing.T) {
	c := New(0.2, 0.3, 0.4)

	assertAlmost(t, c.Mul(2), New(0.4, 0.6, 0.8))
}

func TestHadamardProduct(t *testing.T) {
	c1 := New(1, 0.2, 0.4)
	c2 := New(0.9, 1, 0.1)

	assertAlmost(t, c1.Hadamard(c2), New(0.9, 0.2, 0.04))
}

func assertAlmost(t *testing.T, c1 Color, c2 Color) {
	assert.True(t, float.Equal(c1.R(), c2.R()), "R values differ: c1.R=%v, c2.R=%v", c1.R(), c2.R())
	assert.True(t, float.Equal(c1.G(), c2.G()), "G values differ: c1.G=%v, c2.G=%v", c1.G(), c2.G())
	assert.True(t, float.Equal(c1.B(), c2.B()), "B values differ: c1.B=%v, c2.B=%v", c1.B(), c2.B())
}
