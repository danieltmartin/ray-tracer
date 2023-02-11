package primitive

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/test"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRayMissesCylinder(t *testing.T) {
	c := NewInfCylinder()

	examples := []struct {
		origin, direction tuple.Tuple
	}{
		{tuple.NewPoint(1, 0, 0), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(1, 1, 1)},
	}

	for _, e := range examples {
		r := ray.New(e.origin, e.direction)
		xs := c.localIntersects(r)
		assert.Empty(t, xs)
	}
}

func TestRayHitsCylinder(t *testing.T) {
	c := NewInfCylinder()

	examples := []struct {
		origin, direction tuple.Tuple
		t0, t1            float64
	}{
		{tuple.NewPoint(1, 0, -5), tuple.NewVector(0, 0, 1), 5, 5},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1), 4, 6},
		{tuple.NewPoint(0.5, 0, -5), tuple.NewVector(0.1, 1, 1), 6.80798, 7.08872},
	}

	for _, e := range examples {
		r := ray.New(e.origin, e.direction.Norm())
		xs := c.localIntersects(r)
		require.Len(t, xs, 2)
		test.AssertAlmost(t, e.t0, xs[0].distance)
		test.AssertAlmost(t, e.t1, xs[1].distance)
	}
}

func TestNormalVectorOnCylinder(t *testing.T) {
	c := NewInfCylinder()

	examples := []struct {
		point, normal tuple.Tuple
	}{
		{tuple.NewPoint(1, 0, 0), tuple.NewVector(1, 0, 0)},
		{tuple.NewPoint(0, 5, -1), tuple.NewVector(0, 0, -1)},
		{tuple.NewPoint(0, -2, 1), tuple.NewVector(0, 0, 1)},
		{tuple.NewPoint(-1, 1, 0), tuple.NewVector(-1, 0, 0)},
	}

	for _, e := range examples {
		n := c.localNormalAt(e.point)
		assert.Equal(t, e.normal, n)
	}
}

func TestDefaultMinMaxCylinder(t *testing.T) {
	c := NewInfCylinder()
	assert.Equal(t, math.Inf(-1), c.minY)
	assert.Equal(t, math.Inf(1), c.maxY)
}

func TestIntersectingConstrainedCylinder(t *testing.T) {
	c := NewCylinder(1, 2, false)

	examples := []struct {
		origin, direction tuple.Tuple
		count             int
	}{
		{tuple.NewPoint(0, 1.5, 0), tuple.NewVector(0.1, 1, 0), 0},
		{tuple.NewPoint(0, 3, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 2, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 1, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 1.5, -2), tuple.NewVector(0, 0, 1), 2},
	}

	for _, e := range examples {
		r := ray.New(e.origin, e.direction.Norm())
		xs := c.localIntersects(r)
		assert.Len(t, xs, e.count)
	}
}

func TestIntersectingCapsOfClosedCylinder(t *testing.T) {
	c := NewCylinder(1, 2, true)

	examples := []struct {
		origin, direction tuple.Tuple
		count             int
	}{
		{tuple.NewPoint(0, 3, 0), tuple.NewVector(0, -1, 0), 2},
		{tuple.NewPoint(0, 3, -2), tuple.NewVector(0, -1, 2), 2},
		{tuple.NewPoint(0, 4, -2), tuple.NewVector(0, -1, 1), 2},
		{tuple.NewPoint(0, 0, -2), tuple.NewVector(0, 1, 2), 2},
		{tuple.NewPoint(0, -1, -2), tuple.NewVector(0, 1, 1), 2},
	}

	for _, e := range examples {
		r := ray.New(e.origin, e.direction.Norm())
		xs := c.localIntersects(r)
		assert.Len(t, xs, e.count)
	}
}

func TestNormalVectorOnCylinderEndCaps(t *testing.T) {
	c := NewCylinder(1, 2, true)

	examples := []struct {
		point, normal tuple.Tuple
	}{
		{tuple.NewPoint(0, 1, 0), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(0.5, 1, 0), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(0, 1, 0.5), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(0, 2, 0), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0.5, 2, 0), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0, 2, 0.5), tuple.NewVector(0, 1, 0)},
	}

	for _, e := range examples {
		n := c.localNormalAt(e.point)
		assert.Equal(t, e.normal, n)
	}
}

func TestInfiniteCylinderBounds(t *testing.T) {
	c := NewInfCylinder()

	b := c.bounds()

	assert.Equal(t, tuple.NewPoint(-1, math.Inf(-1), -1), b.min)
	assert.Equal(t, tuple.NewPoint(1, math.Inf(1), 1), b.max)
}

func TestConstrainedCylinderBounds(t *testing.T) {
	c := NewCylinder(-5, 15, true)

	b := c.bounds()

	assert.Equal(t, tuple.NewPoint(-1, -5, -1), b.min)
	assert.Equal(t, tuple.NewPoint(1, 15, 1), b.max)
}
