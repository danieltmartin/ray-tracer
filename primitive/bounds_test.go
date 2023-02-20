package primitive

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/test"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestAddPointsToBoundingBox(t *testing.T) {
	b := NewEmptyBoundingBox()
	p1 := tuple.NewPoint(-5, 2, 0)
	p2 := tuple.NewPoint(7, 0, -3)

	b.AddPoint(p1)
	b.AddPoint(p2)

	assert.Equal(t, tuple.NewPoint(-5, 0, -3), b.min)
	assert.Equal(t, tuple.NewPoint(7, 2, 0), b.max)
}

func TestAddBoundingBoxToBoundingBox(t *testing.T) {
	b1 := NewBoundingBox(tuple.NewPoint(-5, -2, 0), tuple.NewPoint(7, 4, 4))
	b2 := NewBoundingBox(tuple.NewPoint(8, -7, -2), tuple.NewPoint(14, 2, 8))

	b1.AddBox(b2)

	assert.Equal(t, tuple.NewPoint(-5, -7, -2), b1.min)
	assert.Equal(t, tuple.NewPoint(14, 4, 8), b1.max)
}

func TestContainsPoint(t *testing.T) {
	b := NewBoundingBox(tuple.NewPoint(5, -2, 0), tuple.NewPoint(11, 4, 7))

	examples := []struct {
		point  tuple.Tuple
		result bool
	}{
		{tuple.NewPoint(5, -2, 0), true},
		{tuple.NewPoint(11, 4, 7), true},
		{tuple.NewPoint(8, 1, 3), true},
		{tuple.NewPoint(3, 0, 3), false},
		{tuple.NewPoint(8, -4, 3), false},
		{tuple.NewPoint(8, 1, -1), false},
		{tuple.NewPoint(13, 1, 3), false},
		{tuple.NewPoint(8, 5, 3), false},
		{tuple.NewPoint(8, 1, 8), false},
	}

	for _, e := range examples {
		assert.Equal(t, e.result, b.ContainsPoint(e.point))
	}
}

func TestContainsBox(t *testing.T) {
	b := NewBoundingBox(tuple.NewPoint(5, -2, 0), tuple.NewPoint(11, 4, 7))

	examples := []struct {
		min, max tuple.Tuple
		result   bool
	}{
		{tuple.NewPoint(5, -2, 0), tuple.NewPoint(11, 4, 7), true},
		{tuple.NewPoint(6, -1, 1), tuple.NewPoint(10, 3, 6), true},
		{tuple.NewPoint(4, -3, -1), tuple.NewPoint(10, 3, 6), false},
		{tuple.NewPoint(6, -1, 1), tuple.NewPoint(12, 5, 8), false},
	}

	for _, e := range examples {
		b2 := NewBoundingBox(e.min, e.max)
		assert.Equal(t, e.result, b.ContainsBox(b2))
	}
}

func TestTransformBox(t *testing.T) {
	b := NewBoundingBox(tuple.NewPoint(-1, -1, -1), tuple.NewPoint(1, 1, 1))
	tform := transform.RotationX(math.Pi / 4).Mul(transform.RotationY(math.Pi / 4))

	b2 := b.Transform(tform)

	test.AssertAlmost(t, tuple.NewPoint(-1.4142, -1.7071, -1.7071), b2.min)
	test.AssertAlmost(t, tuple.NewPoint(1.4142, 1.7071, 1.7071), b2.max)
}

func TestIntersectBoundingBox(t *testing.T) {
	b := NewBoundingBox(tuple.NewPoint(-1, -1, -1), tuple.NewPoint(1, 1, 1))

	examples := []struct {
		origin, direction tuple.Tuple
		result            bool
	}{
		{tuple.NewPoint(5, 0.5, 0), tuple.NewVector(-1, 0, 0), true},
		{tuple.NewPoint(-5, 0.5, 0), tuple.NewVector(1, 0, 0), true},
		{tuple.NewPoint(0.5, 5, 0), tuple.NewVector(0, -1, 0), true},
		{tuple.NewPoint(0.5, -5, 0), tuple.NewVector(0, 1, 0), true},
		{tuple.NewPoint(0.5, 0, 5), tuple.NewVector(0, 0, -1), true},
		{tuple.NewPoint(0.5, 0, -5), tuple.NewVector(0, 0, 1), true},
		{tuple.NewPoint(0, 0.5, 0), tuple.NewVector(0, 0, 1), true},
		{tuple.NewPoint(-2, 0, 0), tuple.NewVector(2, 4, 6), false},
		{tuple.NewPoint(0, -2, 0), tuple.NewVector(6, 2, 4), false},
		{tuple.NewPoint(0, 0, -2), tuple.NewVector(4, 6, 2), false},
		{tuple.NewPoint(2, 0, 2), tuple.NewVector(0, 0, -1), false},
		{tuple.NewPoint(0, 2, 2), tuple.NewVector(0, -1, 0), false},
		{tuple.NewPoint(2, 2, 0), tuple.NewVector(-1, 0, 0), false},
	}

	for _, e := range examples {
		ray := ray.New(e.origin, e.direction.Norm())
		assert.Equal(t, e.result, b.intersects(ray))
	}
}

func TestIntersectNonCubicBoundingBox(t *testing.T) {
	b := NewBoundingBox(tuple.NewPoint(5, -2, 0), tuple.NewPoint(11, 4, 7))

	examples := []struct {
		origin, direction tuple.Tuple
		result            bool
	}{
		{tuple.NewPoint(15, 1, 2), tuple.NewVector(-1, 0, 0), true},
		{tuple.NewPoint(-5, -1, 4), tuple.NewVector(1, 0, 0), true},
		{tuple.NewPoint(7, 6, 5), tuple.NewVector(0, -1, 0), true},
		{tuple.NewPoint(9, -5, 6), tuple.NewVector(0, 1, 0), true},
		{tuple.NewPoint(8, 2, 12), tuple.NewVector(0, 0, -1), true},
		{tuple.NewPoint(6, 0, -5), tuple.NewVector(0, 0, 1), true},
		{tuple.NewPoint(8, 1, 3.5), tuple.NewVector(0, 0, 1), true},
		{tuple.NewPoint(9, -1, -8), tuple.NewVector(2, 4, 6), false},
		{tuple.NewPoint(8, 3, -4), tuple.NewVector(6, 2, 4), false},
		{tuple.NewPoint(9, -1, -2), tuple.NewVector(4, 6, 2), false},
		{tuple.NewPoint(4, 0, 9), tuple.NewVector(0, 0, -1), false},
		{tuple.NewPoint(8, 6, -1), tuple.NewVector(0, -1, 0), false},
		{tuple.NewPoint(12, 5, 4), tuple.NewVector(-1, 0, 0), false},
	}

	for _, e := range examples {
		ray := ray.New(e.origin, e.direction.Norm())
		assert.Equal(t, e.result, b.intersects(ray))
	}
}
