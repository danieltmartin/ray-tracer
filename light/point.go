package light

import (
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type PointLight struct {
	position  tuple.Tuple
	intensity floatcolor.Float64Color
}

func NewPointLight(position tuple.Tuple, intensity floatcolor.Float64Color) PointLight {
	if !position.IsPoint() {
		panic("light position set to non-position")
	}
	return PointLight{position, intensity}
}

func (p *PointLight) Position() tuple.Tuple {
	return p.position
}

func (p *PointLight) Intensity() floatcolor.Float64Color {
	return p.intensity
}
