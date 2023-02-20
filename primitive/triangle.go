package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Triangle struct {
	p1, p2, p3, e1, e2, normal tuple.Tuple
	data
}

func NewTriangle(p1, p2, p3 tuple.Tuple) Triangle {
	e1 := p2.Sub(p1)
	e2 := p3.Sub(p1)
	normal := e2.Cross(e1).Norm()
	return Triangle{p1, p2, p3, e1, e2, normal, newData()}
}

func (t *Triangle) Vertices() (tuple.Tuple, tuple.Tuple, tuple.Tuple) {
	return t.p1, t.p2, t.p3
}

func (t *Triangle) Intersects(worldRay ray.Ray) Intersections {
	return t.worldIntersects(worldRay, t)
}

func (t *Triangle) NormalAt(worldPoint tuple.Tuple, xn Intersection) tuple.Tuple {
	return t.worldNormalAt(worldPoint, xn, t)
}

func (t *Triangle) localNormalAt(_ tuple.Tuple, _ Intersection) tuple.Tuple {
	return t.normal
}

func (t *Triangle) localIntersects(localRay ray.Ray) Intersections {
	return triangleIntersects(t, localRay, t.p1, t.p2, t.e1, t.e2)
}

func (t *Triangle) Bounds() Bounds {
	x1, y1, z1, _ := t.p1.XYZW()
	x2, y2, z2, _ := t.p2.XYZW()
	x3, y3, z3, _ := t.p3.XYZW()

	minX := math.Min(x1, math.Min(x2, x3))
	minY := math.Min(y1, math.Min(y2, y3))
	minZ := math.Min(z1, math.Min(z2, z3))
	maxX := math.Max(x1, math.Max(x2, x3))
	maxY := math.Max(y1, math.Max(y2, y3))
	maxZ := math.Max(z1, math.Max(z2, z3))

	return Bounds{tuple.NewPoint(minX, minY, minZ), tuple.NewPoint(maxX, maxY, maxZ)}
}

func triangleIntersects(triangle Primitive, localRay ray.Ray, p1, p2, e1, e2 tuple.Tuple) Intersections {
	dirCrossE2 := localRay.Direction().Cross(e2)
	det := e1.Dot(dirCrossE2)
	if float.Equal(det, 0) {
		return nil
	}

	f := 1 / det
	p1ToOrigin := localRay.Origin().Sub(p1)
	u := f * p1ToOrigin.Dot(dirCrossE2)
	if u < 0 || u > 1 {
		return nil
	}

	originCrossE1 := p1ToOrigin.Cross(e1)
	v := f * localRay.Direction().Dot(originCrossE1)
	if v < 0 || u+v > 1 {
		return nil
	}

	d := f * e2.Dot(originCrossE1)
	return NewIntersections(NewIntersectionWithUV(d, u, v, triangle))
}
