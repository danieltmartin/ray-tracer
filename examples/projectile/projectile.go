package main

import (
	"os"

	"github.com/danieltmartin/ray-tracer/canvas"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/image/ppm"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type projectile struct {
	position tuple.Tuple
	velocity tuple.Tuple
}

type environment struct {
	gravity tuple.Tuple
	wind    tuple.Tuple
}

func main() {
	p := projectile{
		position: tuple.NewPoint(0, 1, 0),
		velocity: tuple.NewVector(1, 1.8, 0).Norm().Mul(11.25),
	}

	e := environment{
		gravity: tuple.NewVector(0, -0.1, 0),
		wind:    tuple.NewVector(-0.01, 0, 0),
	}

	c := canvas.New(900, 550)
	c.WritePixel(uint(p.position.X), c.Height()-uint(p.position.Y), floatcolor.Red)

	for p.position.Y > 0 {
		p = tick(e, p)
		c.WritePixel(uint(p.position.X), c.Height()-uint(p.position.Y), floatcolor.Red)
	}
	f, err := os.Create("projectile.ppm")
	if err != nil {
		panic(err)
	}
	ppm.Encode(f, c)
}

func tick(e environment, p projectile) projectile {
	newPosition := p.position.Add(p.velocity)
	newVelocity := p.velocity.Add(e.gravity).Add(e.wind)
	return projectile{newPosition, newVelocity}
}
