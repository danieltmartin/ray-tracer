package floatcolor

import (
	"image/color"
	"math"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/lucasb-eyer/go-colorful"
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

func NewFromInt(c uint32) Float64Color {
	return Float64Color{float64(c>>16&0xFF) / 0xFF, float64(c>>8&0xFF) / 0xFF, float64(c&0xFF) / 0xFF}
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
	return uint32(clamp(c.r) * 0xFFFF),
		uint32(clamp(c.g) * 0xFFFF),
		uint32(clamp(c.b) * 0xFFFF),
		0xFFFF
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

// Equals checks if two colors are equal or have a difference within float.Epsilon
func (c Float64Color) Equals(c2 Float64Color) bool {
	return float.Equal(c.r, c2.r) &&
		float.Equal(c.g, c2.g) &&
		float.Equal(c.b, c2.b)
}

func (c Float64Color) AlmostEquals(c2 Float64Color, epsilon float64) bool {
	return float.AlmostEqual(c.r, c2.r, epsilon) &&
		float.AlmostEqual(c.g, c2.g, epsilon) &&
		float.AlmostEqual(c.b, c2.b, epsilon)
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

func Lerp(c1, c2 Float64Color, t float64) Float64Color {
	color1, _ := colorful.MakeColor(c1)
	color2, _ := colorful.MakeColor(c2)

	lerped := color1.BlendHcl(color2, t).Clamped()
	return Float64Color{lerped.R, lerped.G, lerped.B}
}
