package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type SmoothTriangle struct {
	p1, p2, p3, n1, n2, n3, e1, e2 tuple.Tuple
	data
}

func NewSmoothTriangle(p1, p2, p3, n1, n2, n3 tuple.Tuple) SmoothTriangle {
	e1 := p2.Sub(p1)
	e2 := p3.Sub(p1)
	return SmoothTriangle{p1, p2, p3, n1, n2, n3, e1, e2, newData()}
}

func (t *SmoothTriangle) Vertices() (tuple.Tuple, tuple.Tuple, tuple.Tuple) {
	return t.p1, t.p2, t.p3
}

func (t *SmoothTriangle) Normals() (tuple.Tuple, tuple.Tuple, tuple.Tuple) {
	return t.n1, t.n2, t.n3
}

func (t *SmoothTriangle) Intersects(worldRay ray.Ray) Intersections {
	return t.worldIntersects(worldRay, t)
}

func (t *SmoothTriangle) NormalAt(worldPoint tuple.Tuple, xn Intersection) tuple.Tuple {
	return t.worldNormalAt(worldPoint, xn, t)
}

func (t *SmoothTriangle) localNormalAt(localPoint tuple.Tuple, xn Intersection) tuple.Tuple {
	return t.n2.Mul(xn.u).Add(t.n3.Mul(xn.v).Add(t.n1.Mul(1 - xn.u - xn.v)))
}

func (t *SmoothTriangle) localIntersects(localRay ray.Ray) Intersections {
	return triangleIntersects(t, localRay, t.p1, t.p2, t.e1, t.e2)
}

func (t *SmoothTriangle) Bounds() *BoundingBox {
	x1, y1, z1, _ := t.p1.XYZW()
	x2, y2, z2, _ := t.p2.XYZW()
	x3, y3, z3, _ := t.p3.XYZW()

	minX := math.Min(x1, math.Min(x2, x3))
	minY := math.Min(y1, math.Min(y2, y3))
	minZ := math.Min(z1, math.Min(z2, z3))
	maxX := math.Max(x1, math.Max(x2, x3))
	maxY := math.Max(y1, math.Max(y2, y3))
	maxZ := math.Max(z1, math.Max(z2, z3))

	return &BoundingBox{tuple.NewPoint(minX, minY, minZ), tuple.NewPoint(maxX, maxY, maxZ)}
}
