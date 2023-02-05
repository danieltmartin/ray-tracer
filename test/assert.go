package test

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/stretchr/testify/assert"
)

func AssertAlmost(t *testing.T, c1 any, c2 any) {
	switch c1.(type) {
	case floatcolor.Float64Color:
		r1, g1, b1 := c1.(floatcolor.Float64Color).RGB()
		r2, g2, b2 := c2.(floatcolor.Float64Color).RGB()
		assert.True(t, float.AlmostEqual(r1, r2, 0.001), "R values differ: c1.R=%v, c2.R=%v", r1, r2)
		assert.True(t, float.AlmostEqual(g1, g2, 0.001), "G values differ: c1.G=%v, c2.G=%v", g1, g2)
		assert.True(t, float.AlmostEqual(b1, b2, 0.001), "B values differ: c1.B=%v, c2.B=%v", b1, b2)
	case float64:
		assert.True(t, float.AlmostEqual(c1.(float64), c2.(float64), 0.001), "values differ: c1=%v, c2=%v", c1, c2)
	default:
		panic("unhandled type")
	}
}
