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

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)

	hex := hexagon()
	hex.SetMaterial(material.Default.WithColor(floatcolor.NewFromInt(0xafba3c)))
	hex.SetTransform(transform.Identity().RotationX(-math.Pi/5).Translation(0, 1, 0).Matrix())

	world := world.New()
	world.AddPrimitives(hex)
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

	f, err := os.Create("hexagon.png")
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

func hexagonCorner() primitive.Primitive {
	edge := primitive.NewSphere()
	edge.SetTransform(transform.Identity().
		Scaling(0.25, 0.25, 0.25).
		Translation(0, 0, -1).
		Matrix())
	return &edge
}

func hexagonEdge() primitive.Primitive {
	edge := primitive.NewCylinder(0, 1, false)
	edge.SetTransform(transform.Identity().
		Scaling(0.25, 1, 0.25).
		RotationZ(-math.Pi/2).
		RotationY(-math.Pi/6).
		Translation(0, 0, -1).
		Matrix())
	return &edge
}

func hexagonSide() primitive.Primitive {
	side := primitive.NewGroup()
	side.Add(hexagonCorner(), hexagonEdge())
	return &side
}

func hexagon() primitive.Primitive {
	hex := primitive.NewGroup()
	for n := 0; n <= 5; n++ {
		side := hexagonSide()
		side.SetTransform(transform.RotationY(float64(n) * math.Pi / 3))
		hex.Add(side)
	}
	return &hex
}
