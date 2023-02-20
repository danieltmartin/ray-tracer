package primitive

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestNormalOfPlaneIsConstantEverywhere(t *testing.T) {
	p := NewPlane()

	n1 := p.localNormalAt(tuple.NewPoint(0, 0, 0))
	n2 := p.localNormalAt(tuple.NewPoint(10, 0, -10))
	n3 := p.localNormalAt(tuple.NewPoint(-5, 0, 150))

	up := tuple.NewVector(0, 1, 0)
	assert.Equal(t, up, n1)
	assert.Equal(t, up, n2)
	assert.Equal(t, up, n3)
}

func TestIntersectRayParallelToPlane(t *testing.T) {
	p := NewPlane()
	r := ray.New(tuple.NewPoint(0, 10, 0), tuple.NewVector(0, 0, 1))

	xs := p.localIntersects(r)

	assert.Empty(t, xs)
}

func TestInsersectWithCoplanarRay(t *testing.T) {
	p := NewPlane()
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))

	xs := p.localIntersects(r)

	assert.Empty(t, xs)
}

func TestRayIntersectingPlaneFromAbove(t *testing.T) {
	p := NewPlane()
	r := ray.New(tuple.NewPoint(0, 1, 0), tuple.NewVector(0, -1, 0))

	xs := p.localIntersects(r)

	assert.Len(t, xs, 1)
	assert.Equal(t, 1.0, xs[0].distance)
	assert.Equal(t, &p, xs[0].object)
}

func TestRayIntersectingPlaneFromBelow(t *testing.T) {
	p := NewPlane()
	r := ray.New(tuple.NewPoint(0, -1, 0), tuple.NewVector(0, 1, 0))

	xs := p.localIntersects(r)

	assert.Len(t, xs, 1)
	assert.Equal(t, 1.0, xs[0].distance)
	assert.Equal(t, &p, xs[0].object)
}

func TestPlaneBounds(t *testing.T) {
	p := NewPlane()

	b := p.Bounds()

	assert.Equal(t, tuple.NewPoint(math.Inf(-1), math.Inf(-1), math.Inf(-1)), b.min)
	assert.Equal(t, tuple.NewPoint(math.Inf(1), math.Inf(1), math.Inf(1)), b.max)
}
