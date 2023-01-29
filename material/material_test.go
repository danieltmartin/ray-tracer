package material

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

var obj dummyObject = dummyObject(matrix.Identity4())

func TestDefaultMaterial(t *testing.T) {
	m := Default

	assert.Equal(t, SolidPattern(floatcolor.White), m.Pattern())
	assert.Equal(t, 0.1, m.Ambient())
	assert.Equal(t, 0.9, m.Diffuse())
	assert.Equal(t, 0.9, m.Specular())
	assert.Equal(t, 200.0, m.Shininess())
	assert.Equal(t, 0.0, m.reflective)
}

func TestLightingEyeBetweenLightAndSurface(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)

	color := m.Lighting(obj, light, position, eyev, normalv, false)

	assert.Equal(t, floatcolor.New(1.9, 1.9, 1.9), color)
}

func TestLightingEyeOffset45Degrees(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, math.Sqrt2/2, -math.Sqrt2/2)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)

	color := m.Lighting(obj, light, position, eyev, normalv, false)

	assert.Equal(t, floatcolor.New(1.0, 1.0, 1.0), color)
}

func TestLightingWithLightOffset45Degrees(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 10, -10), floatcolor.White)

	color := m.Lighting(obj, light, position, eyev, normalv, false)

	assert.True(t, floatcolor.New(0.7364, 0.7364, 0.7364).Equals(color))
}

func TestLightingWithEyeInPathOfReflectionVector(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, -math.Sqrt2/2, -math.Sqrt2/2)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 10, -10), floatcolor.White)

	color := m.Lighting(obj, light, position, eyev, normalv, false)

	assert.True(t, floatcolor.New(1.6364, 1.6364, 1.6364).Equals(color))
}

func TestLightingWithLightBehindSurface(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, 10), floatcolor.White)

	color := m.Lighting(obj, light, position, eyev, normalv, false)

	assert.True(t, floatcolor.New(0.1, 0.1, 0.1).Equals(color))
}

func TestLightingWithSurfaceInShadow(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)
	inShadow := true

	color := m.Lighting(obj, light, position, eyev, normalv, inShadow)

	assert.Equal(t, floatcolor.New(0.1, 0.1, 0.1), color)
}

func TestLightingWithStripePattern(t *testing.T) {
	m := Default.
		WithPattern(NewStripePattern(floatcolor.White, floatcolor.Black)).
		WithAmbient(1).
		WithDiffuse(0).
		WithSpecular(0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)

	color1 := m.Lighting(obj, light, tuple.NewPoint(0.9, 0, 0), eyev, normalv, false)
	color2 := m.Lighting(obj, light, tuple.NewPoint(1.1, 0, 0), eyev, normalv, false)

	assert.Equal(t, floatcolor.New(1, 1, 1), color1)
	assert.Equal(t, floatcolor.New(0, 0, 0), color2)
}
