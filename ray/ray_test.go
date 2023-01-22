package ray

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestCreateRay(t *testing.T) {
	origin := tuple.NewPoint(1, 2, 3)
	direction := tuple.NewVector(4, 5, 6)

	r := New(origin, direction)

	assert.Equal(t, origin, r.Origin())
	assert.Equal(t, direction, r.Direction())
}

func TestPosition(t *testing.T) {
	r := New(tuple.NewPoint(2, 3, 4), tuple.NewVector(1, 0, 0))

	assert.Equal(t, tuple.NewPoint(2, 3, 4), r.Position(0))
	assert.Equal(t, tuple.NewPoint(3, 3, 4), r.Position(1))
	assert.Equal(t, tuple.NewPoint(1, 3, 4), r.Position(-1))
	assert.Equal(t, tuple.NewPoint(4.5, 3, 4), r.Position(2.5))
}

func TestTranslatingRay(t *testing.T) {
	r := New(tuple.NewPoint(1, 2, 3), tuple.NewVector(0, 1, 0))
	m := transform.Translation(3, 4, 5)

	r2 := r.Transform(m)

	assert.Equal(t, tuple.NewPoint(4, 6, 8), r2.Origin())
	assert.Equal(t, tuple.NewVector(0, 1, 0), r2.Direction())
}

func TestScalingRay(t *testing.T) {
	r := New(tuple.NewPoint(1, 2, 3), tuple.NewVector(0, 1, 0))
	m := transform.Scaling(2, 3, 4)

	r2 := r.Transform(m)

	assert.Equal(t, tuple.NewPoint(2, 6, 12), r2.Origin())
	assert.Equal(t, tuple.NewVector(0, 3, 0), r2.Direction())
}
