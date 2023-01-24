package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Sphere struct {
	Data
}

func NewSphere() Sphere {
	return Sphere{
		newData(),
	}
}

func (s *Sphere) Intersects(worldRay ray.Ray) Intersections {
	return s.worldIntersects(worldRay, s)
}

func (s *Sphere) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return s.worldNormalAt(worldPoint, s)
}

func (s *Sphere) localIntersects(localRay ray.Ray) Intersections {
	sphereToRay := localRay.Origin().Sub(tuple.NewPoint(0, 0, 0))

	a := localRay.Direction().Dot(localRay.Direction())
	b := 2 * localRay.Direction().Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return nil
	}

	return NewIntersections(
		NewIntersection((-b-math.Sqrt(discriminant))/(2*a), s),
		NewIntersection((-b+math.Sqrt(discriminant))/(2*a), s),
	)
}

func (s *Sphere) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	return localPoint.Sub(tuple.NewPoint(0, 0, 0))
}
