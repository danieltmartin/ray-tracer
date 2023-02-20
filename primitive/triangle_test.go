package primitive

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestConstructTriangle(t *testing.T) {
	p1 := tuple.NewPoint(0, 1, 0)
	p2 := tuple.NewPoint(-1, 0, 0)
	p3 := tuple.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)

	assert.Equal(t, p1, tr.p1)
	assert.Equal(t, p2, tr.p2)
	assert.Equal(t, p3, tr.p3)
	assert.Equal(t, tuple.NewVector(-1, -1, 0), tr.e1)
	assert.Equal(t, tuple.NewVector(1, -1, 0), tr.e2)
	assert.Equal(t, tuple.NewVector(0, 0, -1), tr.normal)
}

func TestNormalOfTriangle(t *testing.T) {
	p1 := tuple.NewPoint(0, 1, 0)
	p2 := tuple.NewPoint(-1, 0, 0)
	p3 := tuple.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)

	n1 := tr.localNormalAt(tuple.NewPoint(0, 0.5, 0))
	n2 := tr.localNormalAt(tuple.NewPoint(-0.5, 0.75, 0))
	n3 := tr.localNormalAt(tuple.NewPoint(0.5, 0.25, 0))

	assert.Equal(t, tr.normal, n1)
	assert.Equal(t, tr.normal, n2)
	assert.Equal(t, tr.normal, n3)
}

func TestIntersectTriangleWithParallelRay(t *testing.T) {
	p1 := tuple.NewPoint(0, 1, 0)
	p2 := tuple.NewPoint(-1, 0, 0)
	p3 := tuple.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)
	r := ray.New(tuple.NewPoint(0, -1, -2), tuple.NewVector(0, 1, 0))

	xs := tr.localIntersects(r)

	assert.Empty(t, xs)
}

func TestIntersectTriangleMissesP1P3Edge(t *testing.T) {
	p1 := tuple.NewPoint(0, 1, 0)
	p2 := tuple.NewPoint(-1, 0, 0)
	p3 := tuple.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)
	r := ray.New(tuple.NewPoint(1, 1, -2), tuple.NewVector(0, 0, 1))

	xs := tr.localIntersects(r)

	assert.Empty(t, xs)
}

func TestIntersectTriangleMissesP1P2Edge(t *testing.T) {
	p1 := tuple.NewPoint(0, 1, 0)
	p2 := tuple.NewPoint(-1, 0, 0)
	p3 := tuple.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)
	r := ray.New(tuple.NewPoint(-1, 1, -2), tuple.NewVector(0, 0, 1))

	xs := tr.localIntersects(r)

	assert.Empty(t, xs)
}

func TestIntersectTriangleMissesP2P3Edge(t *testing.T) {
	p1 := tuple.NewPoint(0, 1, 0)
	p2 := tuple.NewPoint(-1, 0, 0)
	p3 := tuple.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)
	r := ray.New(tuple.NewPoint(0, -1, -2), tuple.NewVector(0, 0, 1))

	xs := tr.localIntersects(r)

	assert.Empty(t, xs)
}

func TestIntersectTriangleHit(t *testing.T) {
	p1 := tuple.NewPoint(0, 1, 0)
	p2 := tuple.NewPoint(-1, 0, 0)
	p3 := tuple.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)
	r := ray.New(tuple.NewPoint(0, 0.5, -2), tuple.NewVector(0, 0, 1))

	xs := tr.localIntersects(r)

	assert.Len(t, xs, 1)
	assert.Equal(t, 2.0, xs[0].distance)
}
