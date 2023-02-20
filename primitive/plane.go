package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Plane struct {
	data
}

func NewPlane() Plane {
	return Plane{newData()}
}

func (p *Plane) NormalAt(worldPoint tuple.Tuple, xn Intersection) tuple.Tuple {
	return p.worldNormalAt(worldPoint, xn, p)
}

func (p *Plane) Intersects(worldRay ray.Ray) Intersections {
	return p.worldIntersects(worldRay, p)
}

func (p *Plane) localIntersects(localRay ray.Ray) Intersections {
	if math.Abs(localRay.Direction().Y) < float.Epsilon {
		return nil
	}
	distance := -localRay.Origin().Y / localRay.Direction().Y
	return NewIntersections(NewIntersection(distance, p))
}

func (p *Plane) localNormalAt(localPoint tuple.Tuple, _ Intersection) tuple.Tuple {
	return tuple.NewVector(0, 1, 0)
}

var planeBounds = NewBoundingBox(
	tuple.NewPoint(math.Inf(-1), 0, math.Inf(-1)),
	tuple.NewPoint(math.Inf(1), 0, math.Inf(1)),
)

func (p *Plane) Bounds() *BoundingBox {
	return planeBounds
}
