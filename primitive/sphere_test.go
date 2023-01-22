package primitive

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRayIntersectsSphereAtTwoPoints(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.Intersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(4.0, &s), xs[0])
	assert.Equal(t, NewIntersection(6.0, &s), xs[1])
}

func TestRayIntersectsSphereAtTangent(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 1, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.Intersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(5.0, &s), xs[0])
	assert.Equal(t, NewIntersection(5.0, &s), xs[1])
}

func TestRayMissesSphere(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 2, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.Intersects(r)

	require.Empty(t, xs)
}

func TestRayInsideSphere(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.Intersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(-1.0, &s), xs[0])
	assert.Equal(t, NewIntersection(1.0, &s), xs[1])
}

func TestSphereBehindRay(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.Intersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(-6.0, &s), xs[0])
	assert.Equal(t, NewIntersection(-4.0, &s), xs[1])
}

func TestSphereDefaultTransformation(t *testing.T) {
	s := NewSphere()

	assert.Equal(t, matrix.Identity4(), s.transform)
}

func TestChangeSphereTransformation(t *testing.T) {
	s := NewSphere()

	transform := transform.Translation(2, 3, 4)
	s.SetTransform(transform)

	assert.Equal(t, transform, s.transform)
}

func TestIntersectingScaledSphere(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	s.SetTransform(transform.Scaling(2, 2, 2))
	xs := s.Intersects(r)

	assert.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(3, &s), xs[0])
	assert.Equal(t, NewIntersection(7, &s), xs[1])
}

func TestIntersectingTranslatedSphere(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	s.SetTransform(transform.Translation(5, 0, 0))
	xs := s.Intersects(r)

	assert.Empty(t, xs)
}

func TestSphereNormalXAxis(t *testing.T) {
	s := NewSphere()

	n := s.NormalAt(tuple.NewPoint(1, 0, 0))

	assert.Equal(t, tuple.NewVector(1, 0, 0), n)
}

func TestSphereNormalYAxis(t *testing.T) {
	s := NewSphere()

	n := s.NormalAt(tuple.NewPoint(0, 1, 0))

	assert.Equal(t, tuple.NewVector(0, 1, 0), n)
}

func TestSphereNormalZAxis(t *testing.T) {
	s := NewSphere()

	n := s.NormalAt(tuple.NewPoint(0, 0, 1))

	assert.Equal(t, tuple.NewVector(0, 0, 1), n)
}

func TestSphereNormalNonAxial(t *testing.T) {
	s := NewSphere()

	p := math.Sqrt(3) / 3
	n := s.NormalAt(tuple.NewPoint(p, p, p))

	assert.True(t, tuple.NewVector(p, p, p).Equals(n))
}

func TestSphereNormalIsNormalized(t *testing.T) {
	s := NewSphere()

	p := math.Sqrt(3) / 3
	n := s.NormalAt(tuple.NewPoint(p, p, p))

	assert.Equal(t, n, n.Norm())
}

func TestSphereNormalTranslated(t *testing.T) {
	s := NewSphere()
	s.SetTransform(transform.Translation(0, 1, 0))

	n := s.NormalAt(tuple.NewPoint(0, 1.70711, -0.70711))

	assert.True(t, tuple.NewVector(0, 0.70711, -0.70711).Equals(n))
}

func TestSphereNormalTransformed(t *testing.T) {
	s := NewSphere()
	s.SetTransform(transform.Identity().RotationZ(math.Pi/5).Scaling(1, 0.5, 1).Matrix())

	n := s.NormalAt(tuple.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))

	assert.True(t, tuple.NewVector(0, 0.97014, -0.24254).Equals(n))
}

func TestSphereDefaultMaterial(t *testing.T) {
	s := NewSphere()

	assert.Equal(t, material.Default, s.Material())
}

func TestSphereAssignMaterial(t *testing.T) {
	s := NewSphere()
	m := material.New(floatcolor.Black, 1, 1, 1, 1)

	s.SetMaterial(m)

	assert.Equal(t, m, s.Material())
}
