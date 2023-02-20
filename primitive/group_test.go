package primitive

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	g := NewGroup()
	assert.Equal(t, matrix.Identity4(), g.transform)
	assert.Empty(t, g.children)
}

func TestAddChildToGroup(t *testing.T) {
	g := NewGroup()
	c := NewCube()
	g.Add(&c)

	assert.Equal(t, &c, g.children[0])
}

func TestAddGroupToItself(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	g := NewGroup()
	g.Add(g)
}

func TestIntersectRayWithEmptyGroup(t *testing.T) {
	g := NewGroup()
	ray := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	xs := g.localIntersects(ray)
	assert.Empty(t, xs)
}

func TestIntersectRayWithNonEmptyGroup(t *testing.T) {
	g := NewGroup()
	s1 := NewSphere()
	s2 := NewSphere()
	s2.SetTransform(transform.Translation(0, 0, -3))
	s3 := NewSphere()
	s3.SetTransform(transform.Translation(5, 0, 0))
	g.Add(&s1, &s2, &s3)

	ray := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	xs := g.localIntersects(ray)

	require.Len(t, xs, 4)
	assert.Equal(t, &s2, xs[0].object)
	assert.Equal(t, &s2, xs[1].object)
	assert.Equal(t, &s1, xs[2].object)
	assert.Equal(t, &s1, xs[3].object)
}

func TestIntersectTransformedGroup(t *testing.T) {
	g := NewGroup()
	g.SetTransform(transform.Scaling(2, 2, 2))
	s := NewSphere()
	s.SetTransform(transform.Translation(5, 0, 0))
	g.Add(&s)

	ray := ray.New(tuple.NewPoint(10, 0, -10), tuple.NewVector(0, 0, 1))
	xs := g.Intersects(ray)

	assert.Len(t, xs, 2)
}

func TestGroupBoundsSingleObject(t *testing.T) {
	g := NewGroup()
	g.SetTransform(transform.Scaling(2, 2, 2))
	s := NewSphere()
	s.SetTransform(transform.Translation(5, 0, 0))
	g.Add(&s)

	b := g.Bounds()

	assert.Equal(t, tuple.NewPoint(4, -1, -1), b.min)
	assert.Equal(t, tuple.NewPoint(6, 1, 1), b.max)
}

func TestGroupBoundsTwoObjects(t *testing.T) {
	g := NewGroup()
	g.SetTransform(transform.Scaling(2, 2, 2))
	s := NewSphere()
	s.SetTransform(transform.Translation(5, 0, 0))
	g.Add(&s)
	c := NewCube()
	c.SetTransform(transform.Translation(-3, -3, -3))
	g.Add(&c)

	b := g.Bounds()

	assert.Equal(t, tuple.NewPoint(-4, -4, -4), b.min)
	assert.Equal(t, tuple.NewPoint(6, 1, 1), b.max)
}

func TestMaterialOfChildIsMaterialOfGroup(t *testing.T) {
	g := NewGroup()
	s := NewSphere()
	g.Add(&s)
	pattern := material.NewCheckerPattern(floatcolor.Black, floatcolor.White)
	g.SetMaterial(material.Default.
		WithPattern(pattern))

	assert.Equal(t, pattern, s.Material().Pattern())
}
