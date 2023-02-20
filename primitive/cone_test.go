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

func TestIntersectConeWithRay(t *testing.T) {
	c := NewInfCone()

	examples := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		t0, t1    float64
	}{
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1), 5, 5},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(1, 1, 1), 8.66025, 8.66025},
		{tuple.NewPoint(1, 1, -5), tuple.NewVector(-0.5, -1, 1), 4.55006, 49.44994},
	}

	for _, e := range examples {
		ray := ray.New(e.origin, e.direction.Norm())
		xs := c.localIntersects(ray)

		require.Len(t, xs, 2)
		test.AssertAlmost(t, e.t0, xs[0].distance)
		test.AssertAlmost(t, e.t1, xs[1].distance)
	}
}

func TestIntersectConeWithRayParallelToOneOfItsHalves(t *testing.T) {
	c := NewInfCone()
	ray := ray.New(tuple.NewPoint(0, 0, -1), tuple.NewVector(0, 1, 1).Norm())

	xs := c.localIntersects(ray)

	require.Len(t, xs, 1)
	test.AssertAlmost(t, 0.35355, xs[0].distance)
}

func TestIntersectConeEndCaps(t *testing.T) {
	c := NewCone(-0.5, 0.5, true)

	examples := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		count     int
	}{
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 1, 0), 0},
		{tuple.NewPoint(0, 0, -0.25), tuple.NewVector(0, 1, 1), 2},
		{tuple.NewPoint(0, 0, -0.25), tuple.NewVector(0, 1, 0), 4},
	}

	for i, e := range examples {
		ray := ray.New(e.origin, e.direction.Norm())
		xs := c.localIntersects(ray)

		require.Len(t, xs, e.count, "expected %v intersections on example %v but got %v", e.count, i, len(xs))
	}
}

func TestNormalVectorOfCone(t *testing.T) {
	c := NewInfCone()

	examples := []struct {
		point, normal tuple.Tuple
	}{
		{tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 0)},
		{tuple.NewPoint(1, 1, 1), tuple.NewVector(1, -math.Sqrt2, 1)},
		{tuple.NewPoint(-1, -1, 0), tuple.NewVector(-1, 1, 0)},
	}

	for _, e := range examples {
		n := c.localNormalAt(e.point)
		assert.Equal(t, e.normal, n)
	}
}

func TestInfiniteConeBounds(t *testing.T) {
	c := NewInfCone()

	b := c.Bounds()

	assert.Equal(t, tuple.NewPoint(-1, math.Inf(-1), -1), b.min)
	assert.Equal(t, tuple.NewPoint(1, math.Inf(1), 1), b.max)
}

func TestConstrainedConeBounds(t *testing.T) {
	c := NewCone(-5, 15, true)

	b := c.Bounds()

	assert.Equal(t, tuple.NewPoint(-1, -5, -1), b.min)
	assert.Equal(t, tuple.NewPoint(1, 15, 1), b.max)
}
