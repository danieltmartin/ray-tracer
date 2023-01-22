package primitive

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/intersect"
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
	assert.Equal(t, intersect.New(4.0, &s), xs[0])
	assert.Equal(t, intersect.New(6.0, &s), xs[1])
}

func TestRayIntersectsSphereAtTangent(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 1, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.Intersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, intersect.New(5.0, &s), xs[0])
	assert.Equal(t, intersect.New(5.0, &s), xs[1])
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
	assert.Equal(t, intersect.New(-1.0, &s), xs[0])
	assert.Equal(t, intersect.New(1.0, &s), xs[1])
}

func TestSphereBehindRay(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	xs := s.Intersects(r)

	require.Len(t, xs, 2)
	assert.Equal(t, intersect.New(-6.0, &s), xs[0])
	assert.Equal(t, intersect.New(-4.0, &s), xs[1])
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
	assert.Equal(t, intersect.New(3, &s), xs[0])
	assert.Equal(t, intersect.New(7, &s), xs[1])
}

func TestIntersectingTranslatedSphere(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := NewSphere()

	s.SetTransform(transform.Translation(5, 0, 0))
	xs := s.Intersects(r)

	assert.Empty(t, xs)
}
