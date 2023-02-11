package main

import (
	"bufio"
	"flag"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/danieltmartin/ray-tracer/camera"
	"github.com/danieltmartin/ray-tracer/floatcolor"
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
			floatcolor.NewFromInt(0xded3d3), floatcolor.Black)).
		WithSpecular(0).
		WithAmbient(0.5).
		WithReflective(0.2),
	)

	wall := primitive.NewPlane()
	wall.SetTransform(transform.Identity().
		RotationX(math.Pi/2).
		RotationY(-math.Pi/4).
		Translation(0, 0, 5).
		Matrix())
	wall.SetMaterial(material.Default.
		WithColor(floatcolor.NewFromInt(0x5d6360)).
		WithDiffuse(0.7).
		WithSpecular(0.8))

	mirror := primitive.NewCube()
	mirror.SetTransform(transform.Identity().
		Scaling(1, 1, 0.05).
		RotationY(-math.Pi/4).
		Translation(-2.4, 1.7, 2).
		Matrix())
	mirror.SetMaterial(material.Default.
		WithDiffuse(0).
		WithAmbient(0).
		WithReflective(1).
		WithRefractiveIndex(1.52).
		WithShininess(500).
		WithSpecular(0.8))

	bigSphere := primitive.NewSphere()
	bigSphere.SetTransform(transform.Identity().RotationX(math.Pi).Translation(-0.5, 1, 0.5).Matrix())
	bigSphere.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0x9f45FF), floatcolor.NewFromInt(0x05020D)).
			WithTransform(transform.RotationY(math.Pi / 8))).
		WithDiffuse(0.7).
		WithReflective(0.2).
		WithSpecular(0.8))

	cube := primitive.NewCube()
	cube.SetTransform(transform.Identity().
		RotationY(math.Pi/5).
		Scaling(0.5, 0.5, 0.5).
		Translation(1.5, 0.5, -1.0).
		Matrix())
	cube.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0x051937), floatcolor.NewFromInt(0xA8EB12))).
		WithDiffuse(0.7).
		WithReflective(0.4).
		WithTransparency(0.2).
		WithRefractiveIndex(2).
		WithSpecular(0.8))

	cyl := primitive.NewCylinder(0, 2, true)
	cyl.SetTransform(transform.Identity().
		RotationX(math.Pi/4).
		Scaling(0.8, 0.8, 0.8).
		Translation(1.5, 1, 2).
		Matrix())
	cyl.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0x051937), floatcolor.NewFromInt(0xA8EB12))).
		WithDiffuse(0.7).
		WithTransparency(0.2).
		WithRefractiveIndex(2.5).
		WithSpecular(0.8))

	cone := primitive.NewCone(0, 1, true)
	cone.SetTransform(transform.Identity().
		RotationZ(math.Pi).
		Translation(0, 1, -0.5).
		Scaling(0.3, 0.3, 0.3).
		Matrix())
	cone.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0x051937), floatcolor.NewFromInt(0xA8EB12)).
			WithTransform(transform.RotationY(math.Pi / 8))).
		WithDiffuse(0.7).
		WithReflective(0.2).
		WithSpecular(0.8))

	smallSphere := primitive.NewSphere()
	smallSphere.SetTransform(transform.Identity().
		Scaling(0.33, 0.33, 0.33).
		Translation(-1.5, 0.33, -0.75).
		Matrix())
	smallSphere.SetMaterial(material.Default.
		WithPattern(material.NewGradientPattern(
			floatcolor.NewFromInt(0x051937), floatcolor.NewFromInt(0xA8EB12)).
			WithTransform(transform.RotationY(math.Pi / 8))).
		WithDiffuse(0.7).
		WithReflective(0.2).
		WithSpecular(0.8))

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)

	group := primitive.NewGroup()
	group.SetTransform(transform.Translation(0, 1, 0))
	group.Add(&floor, &bigSphere, &smallSphere, &cube, &wall, &mirror, &cyl, &cone)

	world := world.New()
	world.AddPrimitives(&group)
	world.AddLights(&light)

	camera := camera.New(1920, 1080, math.Pi/3)
	camera.SetTransform(transform.ViewTransform(
		tuple.NewPoint(0, 1.5, -5),
		tuple.NewPoint(0, 1, 0),
		tuple.NewVector(0, 1, 0),
	))

	start := time.Now()
	image := camera.Render(world)
	duration := time.Since(start)
	log.Printf("Render time: %v\n", duration)

	world.Stats().Log()

	f, err := os.Create("worldcamera.png")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	png.Encode(w, image)
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
