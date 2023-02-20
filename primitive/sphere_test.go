package primitive

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRayIntersectsSphereAtTwoPoints(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.localIntersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(4.0, &s), xs[0])
	assert.Equal(t, NewIntersection(6.0, &s), xs[1])
}

func TestRayIntersectsSphereAtTangent(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 1, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.localIntersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(5.0, &s), xs[0])
	assert.Equal(t, NewIntersection(5.0, &s), xs[1])
}

func TestRayMissesSphere(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 2, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.localIntersects(r)

	require.Empty(t, xs)
}

func TestRayInsideSphere(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.localIntersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(-1.0, &s), xs[0])
	assert.Equal(t, NewIntersection(1.0, &s), xs[1])
}

func TestSphereBehindRay(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.localIntersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, NewIntersection(-6.0, &s), xs[0])
	assert.Equal(t, NewIntersection(-4.0, &s), xs[1])
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

func TestSphereBounds(t *testing.T) {
	s := NewSphere()

	b := s.Bounds()

	assert.Equal(t, tuple.NewPoint(-1, -1, -1), b.min)
	assert.Equal(t, tuple.NewPoint(1, 1, 1), b.max)
}
