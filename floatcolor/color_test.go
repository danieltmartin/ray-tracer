package floatcolor

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/stretchr/testify/assert"
)

func TestNewColor(t *testing.T) {
	c := New(-0.5, 0.4, 1.7)

	r, g, b := c.RGB()
	assert.Equal(t, r, -0.5)
	assert.Equal(t, g, 0.4)
	assert.Equal(t, b, 1.7)
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

func TestRGBAMax(t *testing.T) {
	c := New(1, 1, 1)
	r, g, b, a := c.RGBA()
	assert.Equal(t, uint32(math.MaxUint16), r)
	assert.Equal(t, uint32(math.MaxUint16), g)
	assert.Equal(t, uint32(math.MaxUint16), b)
	assert.Equal(t, uint32(math.MaxUint16), a)
}

func TestRGBAClamp(t *testing.T) {
	c := New(10, 100, 1000000)
	r, g, b, a := c.RGBA()
	assert.Equal(t, uint32(math.MaxUint16), r)
	assert.Equal(t, uint32(math.MaxUint16), g)
	assert.Equal(t, uint32(math.MaxUint16), b)
	assert.Equal(t, uint32(math.MaxUint16), a)
}

func assertAlmost(t *testing.T, c1 Float64Color, c2 Float64Color) {
	r1, g1, b1 := c1.RGB()
	r2, g2, b2 := c2.RGB()
	assert.True(t, float.Equal(r1, r2), "R values differ: c1.R=%v, c2.R=%v", r1, r2)
	assert.True(t, float.Equal(g1, g2), "G values differ: c1.G=%v, c2.G=%v", g1, g2)
	assert.True(t, float.Equal(b1, b2), "B values differ: c1.B=%v, c2.B=%v", b1, b2)
}
