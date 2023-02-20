package primitive

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/test"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestPrimitiveDataDefaults(t *testing.T) {
	d := newData()

	assert.Equal(t, material.Default, d.Material())
	assert.Equal(t, matrix.Identity4(), d.Transform())
}

func TestWorldPointToLocalIdentityTransform(t *testing.T) {
	d := newData()

	local := d.WorldPointToLocal(tuple.NewPoint(1, 2, 3))

	assert.Equal(t, tuple.NewPoint(1, 2, 3), local)
}

func TestWorldPointToLocalTransformed(t *testing.T) {
	d := newData()
	d.SetTransform(transform.Identity().Translation(1, 2, 3).Scaling(2, 2, 2).Matrix())

	local := d.WorldPointToLocal(tuple.NewPoint(1, 1, 1))

	assert.Equal(t, tuple.NewPoint(-0.5, -1.5, -2.5), local)
}

func TestLocalNormalToWorldIdentityTransform(t *testing.T) {
	d := newData()

	world := d.localNormalToWorld(tuple.NewVector(1, 1, 1))

	assert.Equal(t, tuple.NewVector(1, 1, 1).Norm(), world)
}

func TestLocalNormalToWorldTransformed(t *testing.T) {
	d := newData()
	d.SetTransform(transform.Identity().Translation(1, 2, 3).Scaling(2, 1, 2).Matrix())

	world := d.localNormalToWorld(tuple.NewVector(1, 1, 1))

	assert.Equal(t, tuple.NewVector(0.5, 1, 0.5).Norm(), world)
}

func TestLocalNormalToWorldWithParent(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(transform.RotationY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(transform.Scaling(1, 2, 3))
	g1.Add(g2)
	s := NewSphere()
	s.SetTransform(transform.Translation(5, 0, 0))
	g2.Add(&s)

	d := math.Sqrt(3) / 3.0
	p := s.localNormalToWorld(tuple.NewPoint(d, d, d))

	test.AssertAlmost(t, tuple.NewVector(0.2857, 0.4286, -0.8571), p)
}

func TestWorldRayToLocal(t *testing.T) {
	d := newData()
	d.SetTransform(transform.Identity().Translation(1, 2, 3).Scaling(2, 1, 2).Matrix())
	r := ray.New(tuple.NewPoint(1, 1, 1), tuple.NewVector(1, 1, 1))

	local := d.worldRayToLocal(r)

	expected := ray.New(tuple.NewPoint(-0.5, -1, -2.5), tuple.NewVector(0.5, 1, 0.5))
	assert.Equal(t, expected, local)
}

func TestWorldRayToLocalWithParent(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(transform.RotationY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(transform.Scaling(2, 2, 2))
	g1.Add(g2)
	s := NewSphere()
	s.SetTransform(transform.Translation(5, 0, 0))
	g2.Add(&s)

	p := s.WorldPointToLocal(tuple.NewPoint(-2, 0, -10))

	test.AssertAlmost(t, tuple.NewPoint(0, 0, -1), p)
}

func TestWorldNormalOnChildObject(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(transform.RotationY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(transform.Scaling(1, 2, 3))
	g1.Add(g2)
	s := NewSphere()
	s.SetTransform(transform.Translation(5, 0, 0))
	g2.Add(&s)

	p := s.worldNormalAt(tuple.NewPoint(1.7321, 1.1547, -5.5774), Intersection{}, &s)

	test.AssertAlmost(t, tuple.NewVector(0.2857, 0.4286, -0.8571), p)
}
