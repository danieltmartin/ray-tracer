package primitive

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/test"
	"github.com/danieltmartin/ray-tracer/tuple"
)

func TestIntersectionWithSmoothTriangleStoresUV(t *testing.T) {
	tr := NewSmoothTriangle(
		tuple.NewPoint(0, 1, 0),
		tuple.NewPoint(-1, 0, 0),
		tuple.NewPoint(1, 0, 0),
		tuple.NewVector(0, 1, 0),
		tuple.NewVector(-1, 0, 0),
		tuple.NewVector(1, 0, 0),
	)

	r := ray.New(tuple.NewPoint(-0.2, 0.3, -2), tuple.NewVector(0, 0, 1))

	xs := tr.localIntersects(r)

	test.AssertAlmost(t, 0.45, xs[0].u)
	test.AssertAlmost(t, 0.25, xs[0].v)
}

func TestSmoothTriangleUsesUVToInterpolateNormal(t *testing.T) {
	tr := NewSmoothTriangle(
		tuple.NewPoint(0, 1, 0),
		tuple.NewPoint(-1, 0, 0),
		tuple.NewPoint(1, 0, 0),
		tuple.NewVector(0, 1, 0),
		tuple.NewVector(-1, 0, 0),
		tuple.NewVector(1, 0, 0),
	)

	i := NewIntersectionWithUV(1, 0.45, 0.25, &tr)

	n := tr.NormalAt(tuple.NewPoint(0, 0, 0), i)

	test.AssertAlmost(t, tuple.NewVector(-0.5547, 0.83205, 0), n)
}
