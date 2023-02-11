package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Cylinder struct {
	minY, maxY float64
	closed     bool
	data
}

// NewInfCylinder returns a cylinder extending to infinity in the negative and positive y direction.
func NewInfCylinder() Cylinder {
	return Cylinder{math.Inf(-1), math.Inf(1), false, newData()}
}

func NewCylinder(min, max float64, closed bool) Cylinder {
	return Cylinder{min, max, closed, newData()}
}

func (cyl *Cylinder) Intersects(worldRay ray.Ray) Intersections {
	return cyl.worldIntersects(worldRay, cyl)
}

func (cyl *Cylinder) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return cyl.worldNormalAt(worldPoint, cyl)
}

func (cyl *Cylinder) localIntersects(localRay ray.Ray) Intersections {
	direction := localRay.Direction()
	a := direction.X*direction.X + direction.Z*direction.Z

	var xs Intersections

	if !float.Equal(a, 0) {
		origin := localRay.Origin()

		b := 2 * (origin.X*direction.X + origin.Z*direction.Z)
		c := origin.X*origin.X + origin.Z*origin.Z - 1

		discrim := b*b - 4*a*c

		if discrim < 0 {
			return nil
		}

		discrimSqrt := math.Sqrt(discrim)
		t0 := (-b - discrimSqrt) / (2 * a)
		t1 := (-b + discrimSqrt) / (2 * a)

		y0 := origin.Y + t0*direction.Y
		if cyl.minY < y0 && y0 < cyl.maxY {
			xs = append(xs, NewIntersection(t0, cyl))
		}

		y1 := origin.Y + t1*direction.Y
		if cyl.minY < y1 && y1 < cyl.maxY {
			xs = append(xs, NewIntersection(t1, cyl))
		}
	}

	xs = append(xs, cyl.intersectCaps(localRay)...)

	return xs
}

func (cyl *Cylinder) inCap(localRay ray.Ray, t float64) bool {
	x := localRay.Origin().X + t*localRay.Direction().X
	z := localRay.Origin().Z + t*localRay.Direction().Z
	return x*x+z*z <= 1
}

func (cyl *Cylinder) intersectCaps(localRay ray.Ray) Intersections {
	if !cyl.closed || float.Equal(localRay.Direction().Y, 0) {
		return nil
	}

	var xs Intersections

	t := (cyl.minY - localRay.Origin().Y) / localRay.Direction().Y
	if cyl.inCap(localRay, t) {
		xs = append(xs, NewIntersection(t, cyl))
	}

	t = (cyl.maxY - localRay.Origin().Y) / localRay.Direction().Y
	if cyl.inCap(localRay, t) {
		xs = append(xs, NewIntersection(t, cyl))
	}

	return xs
}

func (cyl *Cylinder) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	dist := localPoint.X*localPoint.X + localPoint.Z*localPoint.Z
	if dist < 1 && localPoint.Y >= cyl.maxY-float.Epsilon {
		return tuple.NewVector(0, 1, 0)
	} else if dist < 1 && localPoint.Y <= cyl.minY+float.Epsilon {
		return tuple.NewVector(0, -1, 0)
	}
	return tuple.NewVector(localPoint.X, 0, localPoint.Z)
}
