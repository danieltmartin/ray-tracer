package camera

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/danieltmartin/ray-tracer/world"
	"github.com/stretchr/testify/assert"
)

func TestConstructCamera(t *testing.T) {
	hsize := uint(160)
	vsize := uint(120)
	fieldOfView := math.Pi / 2

	c := New(hsize, vsize, fieldOfView)

	assert.Equal(t, hsize, c.hsize)
	assert.Equal(t, vsize, c.vsize)
	assert.Equal(t, fieldOfView, c.fieldOfView)
	assert.Equal(t, matrix.Identity4(), c.transform)
}

func TestPixelSizeHorizontalCanvas(t *testing.T) {
	c := New(200, 125, math.Pi/2)

	assert.True(t, float.Equal(0.01, c.pixelSize))
}

func TestPixelSizeVerticalCanvas(t *testing.T) {
	c := New(125, 200, math.Pi/2)

	assert.True(t, float.Equal(0.01, c.pixelSize))
}

func TestRayThroughCenterofCanvas(t *testing.T) {
	// Odd numbers so there's an exact center pixel
	c := New(201, 101, math.Pi/2)

	r := c.RayForPixel(100, 50)

	assert.True(t, tuple.NewPoint(0, 0, 0).Equals(r.Origin()))
	assert.True(t, tuple.NewVector(0, 0, -1).Equals(r.Direction()))
}

func TestRayThroughCornerOfCanvas(t *testing.T) {
	c := New(201, 101, math.Pi/2)

	r := c.RayForPixel(0, 0)

	assert.Equal(t, tuple.NewPoint(0, 0, 0), r.Origin())
	assert.True(t, tuple.NewVector(0.66519, 0.33259, -0.66851).Equals(r.Direction()))
}

func TestRayWhenCameraIsTransformed(t *testing.T) {
	c := New(201, 101, math.Pi/2)
	c.SetTransform(transform.Identity().
		Translation(0, -2, 5).
		RotationY(math.Pi / 4).
		Matrix())

	r := c.RayForPixel(100, 50)

	assert.True(t, tuple.NewPoint(0, 2, -5).Equals(r.Origin()))
	assert.True(t, tuple.NewVector(math.Sqrt2/2, 0, -math.Sqrt2/2).Equals(r.Direction()))
}

func TestRenderWorld(t *testing.T) {
	w := testWorld()
	c := New(11, 11, math.Pi/2)
	from := tuple.NewPoint(0, 0, -5)
	to := tuple.NewPoint(0, 0, 0)
	up := tuple.NewVector(0, 1, 0)
	c.SetTransform(transform.ViewTransform(from, to, up))

	image := c.Render(w)

	expected := floatcolor.New(0.38066, 0.47583, 0.2855)
	actual := image.At(5, 5).(floatcolor.Float64Color)
	assert.True(t, expected.Equals(actual))
}

func testWorld() *world.World {
	w := world.New()

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)
	s1 := primitive.NewSphere()
	s1.SetMaterial(material.Default.
		WithColor(floatcolor.New(0.8, 1.0, 0.6)).
		WithDiffuse(0.7).
		WithSpecular(0.2),
	)

	s2 := primitive.NewSphere()
	s2.SetTransform(transform.Scaling(0.5, 0.5, 0.5))

	w.AddLights(&light)
	w.AddPrimitives(&s1, &s2)

	return w
}
