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

const useBVH = true

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

	rand.Seed(0)

	var primitives []primitive.Primitive
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

			primitives = append(primitives, &s)
		}
	}

	rootGroup := primitive.NewGroup()
	if useBVH {
		buildBVH(&rootGroup, primitives)
	} else {
		for _, p := range primitives {
			rootGroup.Add(p)
		}
	}

	world.AddPrimitives(&rootGroup)

	camera := camera.New(1920, 1080, math.Pi/5)
	camera.SetTransform(transform.ViewTransform(
		tuple.NewPoint(0, 50, 0),
		tuple.NewPoint(length*spacing/2, 1, length*spacing/2),
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

func buildBVH(group *primitive.Group, primitives []primitive.Primitive) {
	if len(primitives) <= 4 {
		group.Add(primitives...)
		return
	}
	minX, maxX := xBounds(primitives)
	minZ, maxZ := zBounds(primitives)
	halfX := (maxX + minX) / 2
	halfZ := (maxZ + minZ) / 2

	var partitions [4][]primitive.Primitive
	for _, s := range primitives {
		minBounds := s.Transform().MulTuple(s.Bounds().Min())

		switch {
		case minBounds.X < halfX && minBounds.Z < halfZ:
			partitions[0] = append(partitions[0], s)
		case minBounds.X < halfX && minBounds.Z >= halfZ:
			partitions[1] = append(partitions[1], s)
		case minBounds.X >= halfX && minBounds.Z < halfZ:
			partitions[2] = append(partitions[2], s)
		case minBounds.X >= halfX && minBounds.Z >= halfZ:
			partitions[3] = append(partitions[3], s)
		}
	}

	group1 := primitive.NewGroup()
	group2 := primitive.NewGroup()
	group3 := primitive.NewGroup()
	group4 := primitive.NewGroup()
	buildBVH(&group1, partitions[0])
	buildBVH(&group2, partitions[1])
	buildBVH(&group3, partitions[2])
	buildBVH(&group4, partitions[3])

	group.Add(&group1, &group2, &group3, &group4)
}

func xBounds(spheres []primitive.Primitive) (float64, float64) {
	minX := math.Inf(1)
	maxX := math.Inf(-1)
	for _, s := range spheres {
		minBounds := s.Transform().MulTuple(s.Bounds().Min())
		maxBounds := s.Transform().MulTuple(s.Bounds().Max())
		minX = math.Min(math.Min(minX, minBounds.X), maxBounds.X)
		maxX = math.Max(math.Max(maxX, minBounds.X), maxBounds.X)
	}
	return minX, maxX
}

func zBounds(spheres []primitive.Primitive) (float64, float64) {
	minZ := math.Inf(1)
	maxZ := math.Inf(-1)
	for _, s := range spheres {
		minBounds := s.Transform().MulTuple(s.Bounds().Min())
		mazBounds := s.Transform().MulTuple(s.Bounds().Max())
		minZ = math.Min(math.Min(minZ, minBounds.Z), mazBounds.Z)
		maxZ = math.Max(math.Max(maxZ, minBounds.Z), mazBounds.Z)
	}
	return minZ, maxZ
}
