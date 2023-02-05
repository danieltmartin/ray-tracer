package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Cube struct {
	data
}

func NewCube() Cube {
	return Cube{newData()}
}

func (c *Cube) Intersects(worldRay ray.Ray) Intersections {
	return c.worldIntersects(worldRay, c)
}
func (c *Cube) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return c.worldNormalAt(worldPoint, c)
}

func (c *Cube) localIntersects(localRay ray.Ray) Intersections {
	xtmin, xtmax := checkAxis(localRay.Origin().X, localRay.Direction().X)
	ytmin, ytmax := checkAxis(localRay.Origin().Y, localRay.Direction().Y)
	ztmin, ztmax := checkAxis(localRay.Origin().Z, localRay.Direction().Z)

	tmin := max(xtmin, ytmin, ztmin)
	tmax := min(xtmax, ytmax, ztmax)

	if tmin > tmax {
		return nil
	}

	return NewIntersections(NewIntersection(tmin, c), NewIntersection(tmax, c))
}

func (c *Cube) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	absx := math.Abs(localPoint.X)
	absy := math.Abs(localPoint.Y)
	absz := math.Abs(localPoint.Z)
	maxc := max(absx, absy, absz)
	switch {
	case maxc == absx:
		return tuple.NewVector(localPoint.X, 0, 0)
	case maxc == absy:
		return tuple.NewVector(0, localPoint.Y, 0)
	}
	return tuple.NewVector(0, 0, localPoint.Z)
}

func checkAxis(origin, direction float64) (float64, float64) {
	invDirection := 1 / direction
	tmin := (-1 - origin) * invDirection
	tmax := (1 - origin) * invDirection

	if tmin > tmax {
		return tmax, tmin
	}
	return tmin, tmax
}

func max(a, b, c float64) float64 {
	if a > b && a > c {
		return a
	}
	if b > c {
		return b
	}
	return c
}

func min(a, b, c float64) float64 {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
