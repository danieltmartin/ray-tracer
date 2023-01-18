package color

import "github.com/danieltmartin/ray-tracer/tuple"

type Color struct {
	t tuple.Tuple
}

func New(r, g, b float64) Color {
	return Color{tuple.New(r, g, b, 0)}
}

func (c Color) R() float64 {
	return c.t.X
}

func (c Color) G() float64 {
	return c.t.Y
}

func (c Color) B() float64 {
	return c.t.Z
}

func (c Color) Add(c2 Color) Color {
	return Color{c.t.Add(c2.t)}
}

func (c Color) Sub(c2 Color) Color {
	return Color{c.t.Sub(c2.t)}
}

func (c Color) Mul(v float64) Color {
	return Color{c.t.Mul(v)}
}

func (c Color) Hadamard(c2 Color) Color {
	return Color{c.t.Hadamard(c2.t)}
}
