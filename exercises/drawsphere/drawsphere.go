package main

import (
	"os"

	"github.com/danieltmartin/ray-tracer/canvas"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/image/ppm"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

func main() {
	canvasPixels := 800
	c := canvas.New(canvasPixels, canvasPixels)
	s := primitive.NewSphere()

	rayOrigin := tuple.NewPoint(0, 0, -10)
	wallSize := 7.0
	wallZ := 10.0
	pixelSize := wallSize / float64(canvasPixels)
	half := wallSize / 2

	for y := 0; y < canvasPixels; y++ {
		worldY := half - pixelSize*float64(y)
		for x := 0; x < canvasPixels; x++ {
			worldX := -half + pixelSize*float64(x)

			wallPosition := tuple.NewPoint(worldX, worldY, wallZ)
			r := ray.New(rayOrigin, wallPosition.Sub(rayOrigin).Norm())

			if s.Intersects(r).Hit() != nil {
				c.WritePixel(x, y, floatcolor.Red)
			}
		}
	}

	f, err := os.Create("sphere.ppm")
	if err != nil {
		panic(err)
	}
	ppm.Encode(f, c)
}
