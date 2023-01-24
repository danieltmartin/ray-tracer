package main

import (
	"math"
	"os"

	"github.com/danieltmartin/ray-tracer/camera"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/image/ppm"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/danieltmartin/ray-tracer/world"
)

func main() {
	floor := primitive.NewPlane()
	floor.SetMaterial(floor.Material().
		WithColor(floatcolor.New(1, 0.9, 0.9)).
		WithSpecular(0))

	leftWall := primitive.NewPlane()
	leftWall.SetTransform(transform.Identity().
		RotationX(math.Pi/2).
		RotationY(-math.Pi/4).
		Translation(0, 0, 5).
		Matrix())

	rightWall := primitive.NewPlane()
	rightWall.SetTransform(transform.Identity().
		RotationX(math.Pi/2).
		RotationY(math.Pi/4).
		Translation(0, 0, 5).
		Matrix())

	middleSphere := primitive.NewSphere()
	middleSphere.SetTransform(transform.Translation(-0.5, 1, 0.5))
	middleSphere.SetMaterial(floor.Material().
		WithColor(floatcolor.New(0.1, 1, 0.5)).
		WithDiffuse(0.7).
		WithSpecular(0.3))

	rightSphere := primitive.NewSphere()
	rightSphere.SetTransform(transform.Identity().
		Scaling(0.5, 0.5, 0.5).
		Translation(1.5, 0.5, -0.5).
		Matrix())
	rightSphere.SetMaterial(floor.Material().
		WithColor(floatcolor.New(0.5, 1, 0.1)).
		WithDiffuse(0.7).
		WithSpecular(0.3))

	leftSphere := primitive.NewSphere()
	leftSphere.SetTransform(transform.Identity().
		Scaling(0.33, 0.33, 0.33).
		Translation(-1.5, 0.33, -0.75).
		Matrix())
	leftSphere.SetMaterial(floor.Material().
		WithColor(floatcolor.New(1, 0.8, 0.1)).
		WithDiffuse(0.7).
		WithSpecular(0.3))

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)

	world := world.New()
	world.AddPrimitives(&floor, &leftWall, &rightWall, &middleSphere, &leftSphere, &rightSphere)
	world.AddLights(&light)

	camera := camera.New(400, 300, math.Pi/3)
	camera.SetTransform(transform.ViewTransform(
		tuple.NewPoint(0, 1.5, -5),
		tuple.NewPoint(0, 1, 0),
		tuple.NewVector(0, 1, 0),
	))

	image := camera.Render(world)

	f, err := os.Create("worldcamera.ppm")
	if err != nil {
		panic(err)
	}
	ppm.Encode(f, image)
}
