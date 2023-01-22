package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Sphere struct {
	transform matrix.Matrix
	material  material.Material
}

func NewSphere() Sphere {
	return Sphere{
		matrix.Identity4(),
		material.Default,
	}
}

func (s *Sphere) SetTransform(m matrix.Matrix) {
	s.transform = m
}

func (s *Sphere) SetMaterial(m material.Material) {
	s.material = m
}

func (s *Sphere) Material() material.Material {
	return s.material
}

func (s *Sphere) Intersects(r ray.Ray) Intersections {
	r = r.Transform(s.transform.Inverse())
	sphereToRay := r.Origin().Sub(tuple.NewPoint(0, 0, 0))

	a := r.Direction().Dot(r.Direction())
	b := 2 * r.Direction().Dot(sphereToRay)
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

func (s *Sphere) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	objectPoint := s.transform.Inverse().MulTuple(worldPoint)
	objectNormal := objectPoint.Sub(tuple.NewPoint(0, 0, 0))
	worldNormal := s.transform.Inverse().Transpose().MulTuple(objectNormal)
	return tuple.New(worldNormal.X, worldNormal.Y, worldNormal.Z, 0).Norm()
}
