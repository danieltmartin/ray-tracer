package primitive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersectionEncapsulatesDistanceAndObject(t *testing.T) {
	o := NewSphere()
	i := NewIntersection(3.5, &o)

	assert.Equal(t, 3.5, i.Distance())
	assert.Equal(t, &o, i.Object())
}

func TestAggregatingIntersections(t *testing.T) {
	o := NewSphere()

	i1 := NewIntersection(1, &o)
	i2 := NewIntersection(2, &o)

	xs := NewIntersections(i1, i2)

	assert.Len(t, xs, 2)
	assert.Equal(t, i1, xs[0])
	assert.Equal(t, i2, xs[1])
}

func TestHitAllPositiveDistance(t *testing.T) {
	o := NewSphere()

	i1 := NewIntersection(1, &o)
	i2 := NewIntersection(2, &o)
	xs := NewIntersections(i2, i1)

	i := xs.Hit()

	assert.Equal(t, &i1, i)
}

func TestHitSomeNegative(t *testing.T) {
	o := NewSphere()

	i1 := NewIntersection(-1, &o)
	i2 := NewIntersection(1, &o)
	xs := NewIntersections(i2, i1)

	i := xs.Hit()

	assert.Equal(t, &i2, i)
}

func TestHitAllNegative(t *testing.T) {
	o := NewSphere()

	i1 := NewIntersection(-2, &o)
	i2 := NewIntersection(-1, &o)
	xs := NewIntersections(i2, i1)

	i := xs.Hit()

	assert.Nil(t, i)
}

func TestHitIsAlwaysLowestNonNegativeDistance(t *testing.T) {
	o := NewSphere()

	i1 := NewIntersection(5, &o)
	i2 := NewIntersection(7, &o)
	i3 := NewIntersection(-7, &o)
	i4 := NewIntersection(2, &o)
	xs := NewIntersections(i1, i2, i3, i4)

	i := xs.Hit()

	assert.Equal(t, &i4, i)
}
