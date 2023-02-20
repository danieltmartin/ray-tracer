package primitive

import (
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

func (t *Triangle) Intersects(worldRay ray.Ray) Intersections {
	return t.worldIntersects(worldRay, t)
}

func (t *Triangle) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return t.worldNormalAt(worldPoint, t)
}

func (t *Triangle) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	return t.normal
}

func (t *Triangle) localIntersects(localRay ray.Ray) Intersections {
	dirCrossE2 := localRay.Direction().Cross(t.e2)
	det := t.e1.Dot(dirCrossE2)
	if float.Equal(det, 0) {
		return nil
	}

	f := 1 / det
	p1ToOrigin := localRay.Origin().Sub(t.p1)
	u := f * p1ToOrigin.Dot(dirCrossE2)
	if u < 0 || u > 1 {
		return nil
	}

	originCrossE1 := p1ToOrigin.Cross(t.e1)
	v := f * localRay.Direction().Dot(originCrossE1)
	if v < 0 || u+v > 1 {
		return nil
	}

	d := f * t.e2.Dot(originCrossE1)
	return NewIntersections(NewIntersection(d, t))
}

func (t *Triangle) Bounds() Bounds {
	return Bounds{tuple.NewPoint(0, 0, 0), tuple.NewPoint(0, 0, 0)}
}
