package main

import (
	"os"

	"github.com/danieltmartin/ray-tracer/canvas"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/image/ppm"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

func main() {
	canvasPixels := 800
	c := canvas.New(canvasPixels, canvasPixels)
	s := primitive.NewSphere()

	s.SetMaterial(material.Default.
		WithColor(floatcolor.New(1, 0.2, 1)),
	)

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)

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

			hit := s.Intersects(r).Hit()
			if hit != nil {
				point := r.Position(hit.Distance())
				normal := hit.Object().NormalAt(point)
				eye := r.Direction().Neg()
				color := hit.Object().Material().Lighting(light, point, eye, normal)
				c.WritePixel(x, y, color)
			}
		}
	}

	f, err := os.Create("sphere.ppm")
	if err != nil {
		panic(err)
	}
	ppm.Encode(f, c)
}
