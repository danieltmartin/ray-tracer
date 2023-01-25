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

func (p *Plane) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return p.worldNormalAt(worldPoint, p)
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

func (p *Plane) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	return tuple.NewVector(0, 1, 0)
}
