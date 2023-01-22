package intersect

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/stretchr/testify/assert"
)

type DummyIntersecter bool

func (DummyIntersecter) Intersects(r ray.Ray) Intersections { return nil }

func TestIntersectionEncapsulatesDistanceAndObject(t *testing.T) {
	var o DummyIntersecter
	i := New(3.5, o)

	assert.Equal(t, 3.5, i.Distance())
	assert.Equal(t, o, i.Object())
}

func TestAggregatingIntersections(t *testing.T) {
	var o DummyIntersecter

	i1 := New(1, o)
	i2 := New(2, o)

	xs := NewIntersections(i1, i2)

	assert.Len(t, xs, 2)
	assert.Equal(t, i1, xs[0])
	assert.Equal(t, i2, xs[1])
}

func TestHitAllPositiveDistance(t *testing.T) {
	var o DummyIntersecter

	i1 := New(1, o)
	i2 := New(2, o)
	xs := NewIntersections(i2, i1)

	i := xs.Hit()

	assert.Equal(t, &i1, i)
}

func TestHitSomeNegative(t *testing.T) {
	var o DummyIntersecter

	i1 := New(-1, o)
	i2 := New(1, o)
	xs := NewIntersections(i2, i1)

	i := xs.Hit()

	assert.Equal(t, &i2, i)
}

func TestHitAllNegative(t *testing.T) {
	var o DummyIntersecter

	i1 := New(-2, o)
	i2 := New(-1, o)
	xs := NewIntersections(i2, i1)

	i := xs.Hit()

	assert.Nil(t, i)
}

func TestHitIsAlwaysLowestNonNegativeDistance(t *testing.T) {
	var o DummyIntersecter

	i1 := New(5, o)
	i2 := New(7, o)
	i3 := New(-7, o)
	i4 := New(2, o)
	xs := NewIntersections(i1, i2, i3, i4)

	i := xs.Hit()

	assert.Equal(t, &i4, i)
}
