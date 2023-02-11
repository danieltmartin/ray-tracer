package test

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func AssertAlmost(t *testing.T, c1 any, c2 any) {
	epsilon := 0.0001
	switch c1.(type) {
	case floatcolor.Float64Color:
		r1, g1, b1 := c1.(floatcolor.Float64Color).RGB()
		r2, g2, b2 := c2.(floatcolor.Float64Color).RGB()
		assert.True(t, float.AlmostEqual(r1, r2, epsilon), "R values differ: c1.R=%v, c2.R=%v", r1, r2)
		assert.True(t, float.AlmostEqual(g1, g2, epsilon), "G values differ: c1.G=%v, c2.G=%v", g1, g2)
		assert.True(t, float.AlmostEqual(b1, b2, epsilon), "B values differ: c1.B=%v, c2.B=%v", b1, b2)
	case tuple.Tuple:
		x1, y1, z1, w1 := c1.(tuple.Tuple).XYZW()
		x2, y2, z2, w2 := c2.(tuple.Tuple).XYZW()
		assert.True(t, float.AlmostEqual(x1, x2, epsilon), "X values differ: c1.X=%v, c2.X=%v", x1, x2)
		assert.True(t, float.AlmostEqual(y1, y2, epsilon), "Y values differ: c1.Y=%v, c2.Y=%v", y1, y2)
		assert.True(t, float.AlmostEqual(z1, z2, epsilon), "Z values differ: c1.Z=%v, c2.Z=%v", z1, z2)
		assert.True(t, float.AlmostEqual(w1, w2, epsilon), "W values differ: c1.W=%v, c2.W=%v", w1, w2)
	case float64:
		assert.True(t, float.AlmostEqual(c1.(float64), c2.(float64), epsilon), "values differ: c1=%v, c2=%v", c1, c2)
	default:
		panic("unhandled type")
	}
}
