package material

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMaterial(t *testing.T) {
	m := Default

	assert.Equal(t, floatcolor.White, m.Color())
	assert.Equal(t, 0.1, m.Ambient())
	assert.Equal(t, 0.9, m.Diffuse())
	assert.Equal(t, 0.9, m.Specular())
	assert.Equal(t, 200.0, m.Shininess())
}

func TestLightingEyeBetweenLightAndSurface(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)

	color := m.Lighting(light, position, eyev, normalv)

	assert.Equal(t, floatcolor.New(1.9, 1.9, 1.9), color)
}

func TestLightingEyeOffset45Degrees(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, math.Sqrt2/2, -math.Sqrt2/2)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)

	color := m.Lighting(light, position, eyev, normalv)

	assert.Equal(t, floatcolor.New(1.0, 1.0, 1.0), color)
}

func TestLightingWithLightOffset45Degrees(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 10, -10), floatcolor.White)

	color := m.Lighting(light, position, eyev, normalv)

	assert.True(t, floatcolor.New(0.7364, 0.7364, 0.7364).Equals(color))
}

func TestLightingWithEyeInPathOfReflectionVector(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, -math.Sqrt2/2, -math.Sqrt2/2)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 10, -10), floatcolor.White)

	color := m.Lighting(light, position, eyev, normalv)

	assert.True(t, floatcolor.New(1.6364, 1.6364, 1.6364).Equals(color))
}

func TestLightingWithLightBehindSurface(t *testing.T) {
	m := Default
	position := tuple.NewPoint(0, 0, 0)

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	light := light.NewPointLight(tuple.NewPoint(0, 0, 10), floatcolor.White)

	color := m.Lighting(light, position, eyev, normalv)

	assert.True(t, floatcolor.New(0.1, 0.1, 0.1).Equals(color))
}
