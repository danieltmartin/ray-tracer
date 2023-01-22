package light

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestPointLightHasPositionAndIntensity(t *testing.T) {
	intensity := floatcolor.New(1, 1, 1)
	position := tuple.NewPoint(0, 0, 0)

	light := NewPointLight(position, intensity)

	assert.Equal(t, intensity, light.intensity)
	assert.Equal(t, position, light.position)
}
