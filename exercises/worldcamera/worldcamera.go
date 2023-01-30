package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/danieltmartin/ray-tracer/camera"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/image/ppm"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/danieltmartin/ray-tracer/world"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	floor := primitive.NewPlane()
	floor.SetMaterial(floor.Material().
		WithPattern(material.NewCheckerPattern(
			floatcolor.NewFromInt(0xded3d3), floatcolor.NewFromInt(0x87a9cc))).
		WithSpecular(0).
		WithAmbient(0.5).
		WithReflective(0.2),
	)

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
	middleSphere.SetTransform(transform.Identity().RotationX(math.Pi).Translation(-0.5, 1, 0.5).Matrix())
	middleSphere.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0x9f45FF), floatcolor.NewFromInt(0x05020D)).
			WithTransform(transform.RotationY(math.Pi / 8))).
		WithDiffuse(0.7).
		WithReflective(0.2).
		WithSpecular(0.8))

	rightSphere := primitive.NewSphere()
	rightSphere.SetTransform(transform.Identity().
		Scaling(0.5, 0.5, 0.5).
		Translation(1.5, 0.5, -0.5).
		Matrix())
	rightSphere.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0xC0451A), floatcolor.NewFromInt(0xF0B7B7)).
			WithTransform(transform.RotationY(math.Pi / 8))).
		WithDiffuse(0.7).
		WithReflective(0.2).
		WithSpecular(0.8))

	leftSphere := primitive.NewSphere()
	leftSphere.SetTransform(transform.Identity().
		Scaling(0.33, 0.33, 0.33).
		Translation(-1.5, 0.33, -0.75).
		Matrix())
	leftSphere.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0x051937), floatcolor.NewFromInt(0xA8EB12)).
			WithTransform(transform.RotationY(math.Pi / 8))).
		WithDiffuse(0.7).
		WithReflective(0.2).
		WithSpecular(0.8))

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)

	world := world.New()
	world.AddPrimitives(&floor, &middleSphere, &leftSphere, &rightSphere)
	world.AddLights(&light)

	camera := camera.New(1920, 1080, math.Pi/3)
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
	w := bufio.NewWriter(f)
	ppm.Encode(w, image)
	err = w.Flush()
	if err != nil {
		panic(err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
