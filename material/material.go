package material

import (
	"math"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/tuple"
)

var Default = New(SolidPattern(floatcolor.White), 0.1, 0.9, 0.9, 200.0)

type Material struct {
	pattern                               Pattern
	ambient, diffuse, specular, shininess float64
}

func New(pattern Pattern, ambient, diffuse, specular, shininess float64) Material {
	return Material{pattern, ambient, diffuse, specular, shininess}
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

func (m Material) Pattern() Pattern {
	return m.pattern
}

func (m Material) WithColor(color floatcolor.Float64Color) Material {
	c := m.copy()
	c.pattern = SolidPattern(color)
	return c
}

func (m Material) WithAmbient(a float64) Material {
	c := m.copy()
	c.ambient = a
	return c
}

func (m Material) WithDiffuse(d float64) Material {
	c := m.copy()
	c.diffuse = d
	return c
}

func (m Material) WithSpecular(s float64) Material {
	c := m.copy()
	c.specular = s
	return c
}

func (m Material) WithShininess(s float64) Material {
	c := m.copy()
	c.shininess = s
	return c
}

func (m Material) WithPattern(p Pattern) Material {
	c := m.copy()
	c.pattern = p
	return c
}

func (m Material) copy() Material {
	return Material{m.pattern, m.ambient, m.diffuse, m.specular, m.shininess}
}

func (m Material) Lighting(
	object Object,
	light light.PointLight,
	position tuple.Tuple,
	eyev tuple.Tuple,
	normalv tuple.Tuple,
	inShadow bool,
) floatcolor.Float64Color {
	effectiveColor := m.pattern.colorAtObject(object, position).Hadamard(light.Intensity())
	ambient := effectiveColor.Mul(m.ambient)
	if inShadow {
		return ambient
	}

	// Direction to the light source
	lightv := light.Position().Sub(position).Norm()

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
