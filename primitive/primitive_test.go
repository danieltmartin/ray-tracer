package primitive

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
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

	local := d.worldPointToLocal(tuple.NewPoint(1, 2, 3))

	assert.Equal(t, tuple.NewPoint(1, 2, 3), local)
}

func TestWorldPointToLocalTransformed(t *testing.T) {
	d := newData()
	d.SetTransform(transform.Identity().Translation(1, 2, 3).Scaling(2, 2, 2).Matrix())

	local := d.worldPointToLocal(tuple.NewPoint(1, 1, 1))

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

func TestWorldRayToLocal(t *testing.T) {
	d := newData()
	d.SetTransform(transform.Identity().Translation(1, 2, 3).Scaling(2, 1, 2).Matrix())
	r := ray.New(tuple.NewPoint(1, 1, 1), tuple.NewVector(1, 1, 1))

	local := d.worldRayToLocal(r)

	expected := ray.New(tuple.NewPoint(-0.5, -1, -2.5), tuple.NewVector(0.5, 1, 0.5))
	assert.Equal(t, expected, local)
}
