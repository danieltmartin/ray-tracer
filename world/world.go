package world

import (
	"sort"
	"sync"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type ID uint64

type World struct {
	nextID     ID
	idMutex    sync.Mutex
	primitives []primitive.Primitive
	lights     []*light.PointLight
}

func New() *World {
	return &World{}
}

func (w *World) AddPrimitives(p ...primitive.Primitive) {
	w.primitives = append(w.primitives, p...)
}

func (w *World) AddLights(l ...*light.PointLight) {
	w.lights = append(w.lights, l...)
}

func (w *World) Primitives() []primitive.Primitive {
	return w.primitives
}

func (w *World) Lights() []*light.PointLight {
	return w.lights
}

func (w *World) Intersect(r ray.Ray) primitive.Intersections {
	var allIntersections primitive.Intersections
	for _, p := range w.primitives {
		allIntersections = append(allIntersections, p.Intersects(r)...)
	}

	sort.Slice(allIntersections, func(i, j int) bool {
		return allIntersections[i].Distance() < allIntersections[j].Distance()
	})

	return allIntersections
}

func (w *World) ColorAt(ray ray.Ray) floatcolor.Float64Color {
	hit := w.Intersect(ray).Hit()
	if hit == nil {
		return floatcolor.Black
	}
	hc := prepareHitComputations(*hit, ray)
	return w.shadeHit(hc)
}

func (w *World) shadeHit(hc hitComputations) floatcolor.Float64Color {
	color := floatcolor.Black
	for _, light := range w.lights {
		if light == nil {
			continue
		}
		inShadow := w.isShadowed(hc.overPoint, light.Position())
		hitColor := hc.object.Material().Lighting(hc.object, *light, hc.hitPoint, hc.eyev, hc.normalv, inShadow)
		color = color.Add(hitColor)
	}
	return color
}

func (w *World) isShadowed(p tuple.Tuple, lightPosition tuple.Tuple) bool {
	pointToLight := lightPosition.Sub(p)
	lightDistance := pointToLight.Mag()
	hit := w.Intersect(ray.New(p, pointToLight.Norm())).Hit()
	return hit != nil && hit.Distance() < lightDistance
}

type hitComputations struct {
	distance  float64
	object    primitive.Primitive
	hitPoint  tuple.Tuple
	overPoint tuple.Tuple // Adjusted in normalv direction slightly for floating point precision sensitive calculations
	eyev      tuple.Tuple
	normalv   tuple.Tuple
	inside    bool
}

func prepareHitComputations(intersection primitive.Intersection, ray ray.Ray) hitComputations {
	var hc hitComputations

	hc.distance = intersection.Distance()
	hc.object = intersection.Object()
	hc.hitPoint = ray.Position(intersection.Distance())
	hc.eyev = ray.Direction().Neg()
	hc.normalv = hc.object.NormalAt(hc.hitPoint)

	if hc.normalv.Dot(hc.eyev) < 0 {
		hc.inside = true
		hc.normalv = hc.normalv.Neg()
	}

	hc.overPoint = hc.hitPoint.Add(hc.normalv.Mul(float.Epsilon))

	return hc
}

func (w *World) NextID() ID {
	w.idMutex.Lock()
	defer w.idMutex.Unlock()
	id := w.nextID
	w.nextID += 1
	return id
}
