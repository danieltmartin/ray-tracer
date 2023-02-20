package primitive

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRayIntersectsCube(t *testing.T) {
	c := NewCube()

	examples := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		t1, t2    float64
	}{
		{tuple.NewPoint(5, 0.5, 0), tuple.NewVector(-1, 0, 0), 4, 6}, // +x
		{tuple.NewPoint(-5, 0.5, 0), tuple.NewVector(1, 0, 0), 4, 6}, // -x
		{tuple.NewPoint(0.5, 5, 0), tuple.NewVector(0, -1, 0), 4, 6}, // +y
		{tuple.NewPoint(0.5, -5, 0), tuple.NewVector(0, 1, 0), 4, 6}, // -y
		{tuple.NewPoint(0.5, 0, 5), tuple.NewVector(0, 0, -1), 4, 6}, // +z
		{tuple.NewPoint(0.5, 0, -5), tuple.NewVector(0, 0, 1), 4, 6}, // -z
		{tuple.NewPoint(0, 0.5, 0), tuple.NewVector(0, 0, 1), -1, 1}, // inside
	}

	for _, e := range examples {
		r := ray.New(e.origin, e.direction)
		xs := c.localIntersects(r)

		require.Len(t, xs, 2)
		assert.Equal(t, e.t1, xs[0].distance)
		assert.Equal(t, e.t2, xs[1].distance)
	}
}

func TestRayMissesCube(t *testing.T) {
	c := NewCube()

	examples := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
	}{
		{tuple.NewPoint(-2, 0, 0), tuple.NewVector(0.2673, 0.5345, 0.8018)},
		{tuple.NewPoint(0, -2, 0), tuple.NewVector(0.8018, 0.2673, 0.5345)},
		{tuple.NewPoint(0, 0, -2), tuple.NewVector(0.5345, 0.8018, 0.2673)},
		{tuple.NewPoint(2, 0, 2), tuple.NewVector(0, 0, -1)},
		{tuple.NewPoint(0, 2, 2), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(2, 2, 0), tuple.NewVector(-1, 0, 0)},
	}

	for _, e := range examples {
		r := ray.New(e.origin, e.direction)
		xs := c.localIntersects(r)
		assert.Empty(t, xs)
	}
}

func TestNormalOfCube(t *testing.T) {
	c := NewCube()

	examples := []struct {
		point, normal tuple.Tuple
	}{
		{tuple.NewPoint(1, 0.5, -0.8), tuple.NewVector(1, 0, 0)},
		{tuple.NewPoint(-1, -0.2, 0.9), tuple.NewVector(-1, 0, 0)},
		{tuple.NewPoint(-0.4, 1, -0.1), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0.3, -1, -0.7), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(-0.6, 0.3, 1), tuple.NewVector(0, 0, 1)},
		{tuple.NewPoint(0.4, 0.4, -1), tuple.NewVector(0, 0, -1)},
		{tuple.NewPoint(1, 1, 1), tuple.NewVector(1, 0, 0)},
		{tuple.NewPoint(-1, -1, -1), tuple.NewVector(-1, 0, 0)},
	}

	for _, e := range examples {
		normal := c.localNormalAt(e.point, Intersection{})
		assert.Equal(t, e.normal, normal)
	}
}

func TestCubeBounds(t *testing.T) {
	c := NewCube()

	b := c.Bounds()

	assert.Equal(t, tuple.NewPoint(-1, -1, -1), b.min)
	assert.Equal(t, tuple.NewPoint(1, 1, 1), b.max)
}
