package main

import (
	"math"
	"os"

	"github.com/danieltmartin/ray-tracer/canvas"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/image/ppm"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
)

func main() {
	c := canvas.New(800, 800)

	points := 12
	top := tuple.NewPoint(0, 1, 0)
	for n := 0; n < points; n++ {
		point := transform.Identity().
			RotationZ(float64(n)*2*math.Pi/float64(points)).
			Scaling(350, 350, 0).
			Translation(400, 400, 0).
			Matrix().
			MulTuple(top)
		c.WritePixel(uint(point.X), c.Height()-uint(point.Y), floatcolor.White)
	}

	f, err := os.Create("clock.ppm")
	if err != nil {
		panic(err)
	}
	ppm.Encode(f, c)
}
