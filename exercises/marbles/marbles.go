package main

import (
	"bufio"
	"flag"
	"image/png"
	"log"
	"math"
	"math/rand"
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
		WithReflective(0.0),
	)

	light := light.NewPointLight(tuple.NewPoint(10, 30, -10), floatcolor.White)

	world := world.New()
	world.AddLights(&light)
	world.AddPrimitives(&floor)

	group := primitive.NewGroup()

	rand.Seed(0)

	length := 16.0
	spacing := 3.0
	for x := float64(0); x < length; x++ {
		for z := float64(0); z < length; z++ {
			scale := rand.Float64() + 0.5
			xTrans := x*spacing + spacing*rand.Float64()
			zTrans := z*spacing + spacing*rand.Float64()
			color1 := uint32(rand.Int31())
			color2 := uint32(rand.Int31())
			s := primitive.NewSphere()
			s.SetTransform(transform.Identity().
				Scaling(scale, scale, scale).
				Translation(xTrans, 1+scale/2, zTrans).
				Matrix())
			s.SetMaterial(material.Default.
				WithPattern(material.NewGradientPattern(
					floatcolor.NewFromInt(color1), floatcolor.NewFromInt(color2)).
					WithTransform(transform.RotationY(math.Pi / 8))).
				WithDiffuse(0.7).
				WithReflective(0.2).
				WithSpecular(0.8))

			group.Add(&s)
		}
	}

	world.AddPrimitives(&group)

	camera := camera.New(1920, 1080, math.Pi/7)
	camera.SetTransform(transform.ViewTransform(
		tuple.NewPoint(0, 50, 0),
		tuple.NewPoint(length*spacing, 1, length*spacing),
		// tuple.NewPoint(length*spacing+900, 1, length*spacing+900), // not looking at anything
		tuple.NewVector(0, 1, 0),
	))

	start := time.Now()
	image := camera.Render(world)
	duration := time.Since(start)
	log.Printf("Render time: %v\n", duration)

	world.Stats().Log()

	f, err := os.Create("marbles.png")
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
