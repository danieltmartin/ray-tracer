package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/intersect"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Sphere struct {
	transform matrix.Matrix
}

func NewSphere() Sphere {
	return Sphere{
		matrix.Identity4(),
	}
}

func (s *Sphere) SetTransform(m matrix.Matrix) {
	s.transform = m
}

func (s *Sphere) Intersects(r ray.Ray) intersect.Intersections {
	r = r.Transform(s.transform.Inverse())
	sphereToRay := r.Origin().Sub(tuple.NewPoint(0, 0, 0))

	a := r.Direction().Dot(r.Direction())
	b := 2 * r.Direction().Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return nil
	}

	return intersect.NewIntersections(
		intersect.New((-b-math.Sqrt(discriminant))/(2*a), s),
		intersect.New((-b+math.Sqrt(discriminant))/(2*a), s),
	)
}
