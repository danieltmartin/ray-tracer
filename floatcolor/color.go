package floatcolor

import (
	"image/color"
	"math"

	"github.com/danieltmartin/ray-tracer/float"
)

type Float64Color struct {
	r, g, b float64
}

var (
	Black = New(0, 0, 0)
	White = New(1, 1, 1)
	Red   = New(1, 0, 0)
	Green = New(0, 1, 0)
	Blue  = New(0, 0, 1)
)

func New(r, g, b float64) Float64Color {
	return Float64Color{r, g, b}
}

func (c Float64Color) RGBA() (r, g, b, a uint32) {
	clamp := func(v float64) float64 {
		if v > 1 {
			return 1
		}
		if v < 0 {
			return 0
		}
		return v
	}
	return uint32(clamp(c.r) * math.MaxUint16),
		uint32(clamp(c.g) * math.MaxUint16),
		uint32(clamp(c.b) * math.MaxUint16),
		math.MaxUint16
}

func (c Float64Color) RGB() (r, g, b float64) {
	return c.r, c.g, c.b
}

func (c Float64Color) Add(c2 Float64Color) Float64Color {
	return Float64Color{c.r + c2.r, c.g + c2.g, c.b + c2.b}
}

func (c Float64Color) Sub(c2 Float64Color) Float64Color {
	return Float64Color{c.r - c2.r, c.g - c2.g, c.b - c2.b}
}

func (c Float64Color) Mul(v float64) Float64Color {
	return Float64Color{c.r * v, c.g * v, c.b * v}
}

func (c Float64Color) Hadamard(c2 Float64Color) Float64Color {
	return Float64Color{c.r * c2.r, c.g * c2.g, c.b * c2.b}
}

func (c Float64Color) Equals(c2 Float64Color) bool {
	return float.Equal(c.r, c2.r) && float.Equal(c.g, c2.g) && float.Equal(c.b, c2.b)
}

var Float64Model = color.ModelFunc(float64Model)

func float64Model(c color.Color) color.Color {
	if _, ok := c.(Float64Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return New(
		float64(r)/math.MaxUint32,
		float64(g)/math.MaxUint32,
		float64(b)/math.MaxUint32,
	)
}
