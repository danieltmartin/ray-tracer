package primitive

import "github.com/danieltmartin/ray-tracer/ray"

type Intersecter interface {
	Intersects(r ray.Ray) Intersections
}

type Intersection struct {
	distance float64
	object   Primitive
}

func NewIntersection(distance float64, object Primitive) Intersection {
	return Intersection{distance, object}
}

func (i Intersection) Distance() float64 {
	return i.distance
}

func (i Intersection) Object() Primitive {
	return i.object
}

type Intersections []Intersection

func NewIntersections(i ...Intersection) Intersections {
	return i
}

func (i Intersections) Hit() *Intersection {
	var lowestNonNegative *Intersection

	for n := range i {
		if i[n].distance >= 0 && (lowestNonNegative == nil || i[n].distance < lowestNonNegative.distance) {
			lowestNonNegative = &i[n]
		}
	}
	return lowestNonNegative
}
