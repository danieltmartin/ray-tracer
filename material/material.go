package material

import (
	"math"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/tuple"
)

var Default = New(floatcolor.White, 0.1, 0.9, 0.9, 200.0)

type Material struct {
	color                                 floatcolor.Float64Color
	ambient, diffuse, specular, shininess float64
}

func New(color floatcolor.Float64Color, ambient, diffuse, specular, shininess float64) Material {
	return Material{color, ambient, diffuse, specular, shininess}
}

func (m Material) Color() floatcolor.Float64Color {
	return m.color
}

func (m Material) Ambient() float64 {
	return m.ambient
}

func (m Material) Diffuse() float64 {
	return m.diffuse
}

func (m Material) Specular() float64 {
	return m.specular
}

func (m Material) Shininess() float64 {
	return m.shininess
}

func (m Material) WithColor(color floatcolor.Float64Color) Material {
	return New(color, m.ambient, m.diffuse, m.specular, m.shininess)
}

func (m Material) WithAmbient(a float64) Material {
	return New(m.color, a, m.diffuse, m.specular, m.shininess)
}

func (m Material) WithDiffuse(d float64) Material {
	return New(m.color, m.ambient, d, m.specular, m.shininess)
}

func (m Material) WithSpecular(s float64) Material {
	return New(m.color, m.ambient, m.diffuse, s, m.shininess)
}

func (m Material) WithShininess(s float64) Material {
	return New(m.color, m.ambient, m.diffuse, m.specular, s)
}

func (m Material) Lighting(
	light light.PointLight,
	position tuple.Tuple,
	eyev tuple.Tuple,
	normalv tuple.Tuple) floatcolor.Float64Color {
	effectiveColor := m.color.Hadamard(light.Intensity())

	// Direction to the light source
	lightv := light.Position().Sub(position).Norm()

	ambient := effectiveColor.Mul(m.ambient)
	diffuse := floatcolor.Black
	specular := floatcolor.Black

	// Cosine of angle between light and normal. Negative means the light is on the other side of the surface.
	lightDotNormal := lightv.Dot(normalv)
	if lightDotNormal >= 0 {
		diffuse = effectiveColor.Mul(m.diffuse).Mul(lightDotNormal)
		reflectv := lightv.Neg().Reflect(normalv)

		// Cosine of angle between reflection and eye. Negative means the light reflects away from the eye.
		reflectDotEye := reflectv.Dot(eyev)
		if reflectDotEye >= 0 {
			factor := math.Pow(reflectDotEye, m.shininess)
			specular = light.Intensity().Mul(m.specular * factor)
		}
	}

	return ambient.Add(diffuse).Add(specular)
}