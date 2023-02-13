package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Cone struct {
	minY, maxY float64
	closed     bool
	data
}

// NewInfCylinder returns a cone extending to infinity in the negative and positive y direction.
func NewInfCone() Cone {
	return Cone{math.Inf(-1), math.Inf(1), false, newData()}
}

func NewCone(min, max float64, closed bool) Cone {
	return Cone{min, max, closed, newData()}
}

func (co *Cone) Intersects(worldRay ray.Ray) Intersections {
	return co.worldIntersects(worldRay, co)
}

func (co *Cone) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return co.worldNormalAt(worldPoint, co)
}

func (co *Cone) localIntersects(localRay ray.Ray) Intersections {
	direction := localRay.Direction()
	a := direction.X*direction.X - direction.Y*direction.Y + direction.Z*direction.Z

	var xs Intersections

	origin := localRay.Origin()

	b := 2 * (origin.X*direction.X - origin.Y*direction.Y + origin.Z*direction.Z)
	c := origin.X*origin.X - origin.Y*origin.Y + origin.Z*origin.Z
	if float.Equal(a, 0) {
		if !float.Equal(b, 0) {
			xs = append(xs, NewIntersection(-c/(2*b), co))
		}
	} else {
		discrim := b*b - 4*a*c

		if discrim < 0 {
			return nil
		}

		discrimSqrt := math.Sqrt(discrim)
		t0 := (-b - discrimSqrt) / (2 * a)
		t1 := (-b + discrimSqrt) / (2 * a)

		y0 := origin.Y + t0*direction.Y
		if co.minY < y0 && y0 < co.maxY {
			xs = append(xs, NewIntersection(t0, co))
		}

		y1 := origin.Y + t1*direction.Y
		if co.minY < y1 && y1 < co.maxY {
			xs = append(xs, NewIntersection(t1, co))
		}
	}

	xs = append(xs, co.intersectCaps(localRay)...)

	return xs
}

func (co *Cone) inCap(localRay ray.Ray, t float64, capRadius float64) bool {
	x := localRay.Origin().X + t*localRay.Direction().X
	z := localRay.Origin().Z + t*localRay.Direction().Z
	return x*x+z*z <= capRadius
}

func (co *Cone) intersectCaps(localRay ray.Ray) Intersections {
	if !co.closed || float.Equal(localRay.Direction().Y, 0) {
		return nil
	}

	var xs Intersections

	t := (co.minY - localRay.Origin().Y) / localRay.Direction().Y
	if co.inCap(localRay, t, math.Abs(co.minY)) {
		xs = append(xs, NewIntersection(t, co))
	}

	t = (co.maxY - localRay.Origin().Y) / localRay.Direction().Y
	if co.inCap(localRay, t, math.Abs(co.maxY)) {
		xs = append(xs, NewIntersection(t, co))
	}

	return xs
}

func (co *Cone) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	dist := localPoint.X*localPoint.X + localPoint.Z*localPoint.Z
	if dist < 1 && localPoint.Y >= co.maxY-float.Epsilon {
		return tuple.NewVector(0, 1, 0)
	} else if dist < 1 && localPoint.Y <= co.minY+float.Epsilon {
		return tuple.NewVector(0, -1, 0)
	}
	y := math.Sqrt(localPoint.X*localPoint.X + localPoint.Z*localPoint.Z)
	if localPoint.Y > 0 {
		y = -y
	}
	return tuple.NewVector(localPoint.X, y, localPoint.Z)
}

func (c *Cone) Bounds() Bounds {
	return newBounds(
		tuple.NewPoint(-1, c.minY, -1),
		tuple.NewPoint(1, c.maxY, 1))
}
