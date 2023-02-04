package world

import (
	"math"
	"sort"
	"sync"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

var intersectionPool sync.Pool

func init() {
	intersectionPool.New = func() any {
		xns := primitive.Intersections(make(primitive.Intersections, 0))
		return &xns
	}
}

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

func (w *World) intersect(r ray.Ray) primitive.Intersections {
	allIntersections := *intersectionPool.Get().(*primitive.Intersections)
	for _, p := range w.primitives {
		allIntersections = append(allIntersections, p.Intersects(r)...)
	}

	sort.Slice(allIntersections, func(i, j int) bool {
		return allIntersections[i].Distance() < allIntersections[j].Distance()
	})

	return allIntersections
}

func (w *World) ColorAt(ray ray.Ray, remaining int) floatcolor.Float64Color {
	xns := w.intersect(ray)
	hit := xns.Hit()
	if hit == nil {
		return floatcolor.Black
	}
	hc := prepareHitComputations(*hit, ray, xns...)
	reslice := xns[:0]
	intersectionPool.Put(&reslice)
	return w.shadeHit(hc, remaining)
}

func (w *World) shadeHit(hc hitComputations, remaining int) floatcolor.Float64Color {
	surfaceColor := floatcolor.Black
	for _, light := range w.lights {
		if light == nil {
			continue
		}
		inShadow := w.isShadowed(hc.overPoint, light.Position())
		hitColor := hc.object.Material().Lighting(hc.object, *light, hc.overPoint, hc.eyev, hc.normalv, inShadow)
		surfaceColor = surfaceColor.Add(hitColor)
	}
	reflectColor := w.reflectedColor(hc, remaining)
	refractColor := w.refractedColor(hc, remaining)
	if hc.object.Material().Reflective() > 0 && hc.object.Material().Transparency() > 0 {
		reflectance := schlick(hc)
		return surfaceColor.Add(reflectColor.Mul(reflectance)).Add(refractColor.Mul(1 - reflectance))
	}
	return surfaceColor.Add(reflectColor).Add(refractColor)
}

func (w *World) isShadowed(p tuple.Tuple, lightPosition tuple.Tuple) bool {
	pointToLight := lightPosition.Sub(p)
	lightDistance := pointToLight.Mag()
	xns := w.intersect(ray.New(p, pointToLight.Norm()))
	hit := xns.Hit()
	reslice := xns[:0]
	intersectionPool.Put(&reslice)
	return hit != nil && hit.Distance() < lightDistance
}

func (w *World) reflectedColor(hc hitComputations, remaining int) floatcolor.Float64Color {
	if remaining == 0 || hc.object.Material().Reflective() == 0 {
		return floatcolor.Black
	}
	reflectRay := ray.New(hc.overPoint, hc.reflectv)
	return w.ColorAt(reflectRay, remaining-1).Mul(hc.object.Material().Reflective())
}

func (w *World) refractedColor(hc hitComputations, remaining int) floatcolor.Float64Color {
	if remaining == 0 {
		return floatcolor.Black
	}
	if hc.object.Material().Transparency() == 0 {
		return floatcolor.Black
	}

	// Total internal reflection check
	nRatio := hc.n1 / hc.n2
	cosi := hc.eyev.Dot(hc.normalv)
	sin2t := nRatio * nRatio * (1 - cosi*cosi)
	if sin2t > 1 {
		return floatcolor.Black
	}

	cost := math.Sqrt(1.0 - sin2t)
	direction := hc.normalv.Mul(nRatio*cosi - cost).Sub(hc.eyev.Mul(nRatio))
	refractRay := ray.New(hc.underPoint, direction)

	return w.ColorAt(refractRay, remaining-1).Mul(hc.object.Material().Transparency())
}

type hitComputations struct {
	distance   float64
	object     primitive.Primitive
	hitPoint   tuple.Tuple
	overPoint  tuple.Tuple // Adjusted in normalv direction slightly for floating point precision sensitive calculations
	underPoint tuple.Tuple
	eyev       tuple.Tuple
	normalv    tuple.Tuple
	reflectv   tuple.Tuple
	inside     bool
	n1         float64
	n2         float64
}

func prepareHitComputations(
	hit primitive.Intersection,
	ray ray.Ray,
	allIntersections ...primitive.Intersection) hitComputations {
	var hc hitComputations

	hc.distance = hit.Distance()
	hc.object = hit.Object()
	hc.hitPoint = ray.Position(hit.Distance())
	hc.eyev = ray.Direction().Neg()
	hc.normalv = hc.object.NormalAt(hc.hitPoint)
	hc.reflectv = ray.Direction().Reflect(hc.normalv)

	if hc.normalv.Dot(hc.eyev) < 0 {
		hc.inside = true
		hc.normalv = hc.normalv.Neg()
	}

	hc.overPoint = hc.hitPoint.Add(hc.normalv.Mul(float.Epsilon))
	hc.underPoint = hc.hitPoint.Sub(hc.normalv.Mul(float.Epsilon))

	var containers []primitive.Primitive
	for _, x := range allIntersections {
		if x == hit {
			if len(containers) == 0 {
				hc.n1 = 1.0
			} else {
				hc.n1 = containers[len(containers)-1].Material().RefractiveIndex()
			}
		}

		xRemoved := remove(containers, x.Object())
		if xRemoved != nil {
			containers = xRemoved
		} else {
			containers = append(containers, x.Object())
		}

		if x == hit {
			if len(containers) == 0 {
				hc.n2 = 1.0
			} else {
				hc.n2 = containers[len(containers)-1].Material().RefractiveIndex()
			}
			break
		}
	}

	return hc
}

// schlick computes Schlick's approximation of the Fresnel effect.
func schlick(hc hitComputations) float64 {
	cos := hc.eyev.Dot(hc.normalv)

	if hc.n1 > hc.n2 {
		nRatio := hc.n1 / hc.n2
		sin2t := nRatio * nRatio * (1 - cos*cos)
		if sin2t > 1 { // Total Internal Reflection
			return 1
		}

		cos = math.Sqrt(1 - sin2t)
	}

	r0 := (hc.n1 - hc.n2) / (hc.n1 + hc.n2)
	r0 *= r0
	x := 1 - cos
	return r0 + (1-r0)*x*x*x*x*x
}

func remove(containers []primitive.Primitive, o primitive.Primitive) []primitive.Primitive {
	for i := range containers {
		if containers[i] == o {
			return append(containers[:i], containers[i+1:]...)
		}
	}
	return nil
}

func (w *World) NextID() ID {
	w.idMutex.Lock()
	defer w.idMutex.Unlock()
	id := w.nextID
	w.nextID += 1
	return id
}
