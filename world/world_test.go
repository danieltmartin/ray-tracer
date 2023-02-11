package world

import (
	"math"
	"sync"
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/test"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyWorld(t *testing.T) {
	w := New()

	assert.Empty(t, w.Primitives())
	assert.Empty(t, w.Lights())
}

func TestRayIntersectWorld(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))

	xs := w.intersect(r)

	assert.Len(t, xs, 4)
	assert.Equal(t, 4.0, xs[0].Distance())
	assert.Equal(t, 4.5, xs[1].Distance())
	assert.Equal(t, 5.5, xs[2].Distance())
	assert.Equal(t, 6.0, xs[3].Distance())
}

func TestPrecomputingIntersectionState(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := primitive.NewSphere()
	i := primitive.NewIntersection(4, &s)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, &s, hc.object)
	assert.Equal(t, tuple.NewPoint(0, 0, -1), hc.hitPoint)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.eyev)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.normalv)
}

func TestPrecomputingReflectionVector(t *testing.T) {
	s := primitive.NewPlane()
	r := ray.New(tuple.NewPoint(0, 1, -1), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	i := primitive.NewIntersection(math.Sqrt2, &s)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, tuple.NewVector(0, math.Sqrt2/2, math.Sqrt2/2), hc.reflectv)
}

func TestPrecomputingOverpoint(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := primitive.NewSphere()
	s.SetTransform(transform.Translation(0, 0, 1))
	i := primitive.NewIntersection(5, &s)

	hc := prepareHitComputations(i, r)

	assert.Less(t, hc.overPoint.Z, -float.Epsilon/2)
	assert.Greater(t, hc.hitPoint.Z, hc.overPoint.Z)
}

func TestPrecomputingUnderPoint(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := glassSphere()
	s.SetTransform(transform.Translation(0, 0, 1))
	i := primitive.NewIntersection(5, s)
	xs := primitive.NewIntersections(i)

	hc := prepareHitComputations(i, r, xs...)

	assert.Greater(t, hc.underPoint.Z, float.Epsilon/2)
	assert.Less(t, hc.hitPoint.Z, hc.underPoint.Z)
}

func TestPrecomputingEnterAndExitPoint(t *testing.T) {
	a := glassSphere()
	a.SetTransform(transform.Scaling(2, 2, 2))
	a.SetMaterial(a.Material().WithRefractiveIndex(1.5))

	b := glassSphere()
	b.SetTransform(transform.Translation(0, 0, -0.25))
	b.SetMaterial(b.Material().WithRefractiveIndex(2.0))

	c := glassSphere()
	c.SetTransform(transform.Translation(0, 0, 0.25))
	c.SetMaterial(c.Material().WithRefractiveIndex(2.5))

	r := ray.New(tuple.NewPoint(0, 0, -4), tuple.NewVector(0, 0, 1))

	xs := primitive.NewIntersections(
		primitive.NewIntersection(2, a),
		primitive.NewIntersection(2.75, b),
		primitive.NewIntersection(3.25, c),
		primitive.NewIntersection(4.75, b),
		primitive.NewIntersection(5.25, c),
		primitive.NewIntersection(6, a),
	)

	expected := [][]float64{
		{1.0, 1.5},
		{1.5, 2.0},
		{2.0, 2.5},
		{2.5, 2.5},
		{2.5, 1.5},
		{1.5, 1.0},
	}

	for i, x := range xs {
		hc := prepareHitComputations(x, r, xs...)
		assert.Equal(t, expected[i][0], hc.n1)
		assert.Equal(t, expected[i][1], hc.n2)
	}
}

func TestHitOutside(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[0]
	i := primitive.NewIntersection(4, shape)

	hc := prepareHitComputations(i, r)

	assert.False(t, hc.inside)
}

func TestHitInside(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[0]
	i := primitive.NewIntersection(1, shape)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, tuple.NewPoint(0, 0, 1), hc.hitPoint)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.eyev)
	assert.True(t, hc.inside)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.normalv)
}

func TestShadingAnIntersection(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[0]
	i := primitive.NewIntersection(4, shape)

	hc := prepareHitComputations(i, r)
	c := w.shadeHit(hc, 1)

	assert.True(t, floatcolor.New(0.38066, 0.47583, 0.2855).Equals(c))
}

func TestShadingAnIntersectionFromTheInside(t *testing.T) {
	w := testWorld()
	l := light.NewPointLight(tuple.NewPoint(0, 0.25, 0), floatcolor.White)
	w.lights = []*light.PointLight{&l}
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[1]
	i := primitive.NewIntersection(0.5, shape)

	hc := prepareHitComputations(i, r)
	c := w.shadeHit(hc, 1)

	assert.True(t, floatcolor.New(0.90498, 0.90498, 0.90498).Equals(c))
}

func TestShadingAnIntersectionInAShadow(t *testing.T) {
	w := New()
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)
	w.AddLights(&light)
	s1 := primitive.NewSphere()
	s2 := primitive.NewSphere()
	s2.SetTransform(transform.Translation(0, 0, 10))
	w.AddPrimitives(&s1, &s2)
	r := ray.New(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	i := primitive.NewIntersection(4, &s2)

	hc := prepareHitComputations(i, r)
	c := w.shadeHit(hc, 1)

	assert.True(t, floatcolor.New(0.1, 0.1, 0.1).Equals(c))
}

func TestShadingWithTransparentMaterial(t *testing.T) {
	w := testWorld()

	floor := primitive.NewPlane()
	floor.SetTransform(transform.Translation(0, -1, 0))
	floor.SetMaterial(floor.Material().
		WithTransparency(0.5).
		WithRefractiveIndex(1.5),
	)
	w.AddPrimitives(&floor)

	ball := primitive.NewSphere()
	ball.SetTransform(transform.Translation(0, -3.5, -0.5))
	ball.SetMaterial(ball.Material().
		WithColor(floatcolor.Red).
		WithAmbient(0.5),
	)
	w.AddPrimitives(&ball)

	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	x := primitive.NewIntersection(math.Sqrt2, &floor)

	hc := prepareHitComputations(x, r, x)
	color := w.shadeHit(hc, 5)

	test.AssertAlmost(t, floatcolor.New(0.93642, 0.68642, 0.68642), color)
}

func TestShadingWithReflectiveAndTransparentMaterial(t *testing.T) {
	w := testWorld()

	floor := primitive.NewPlane()
	floor.SetTransform(transform.Translation(0, -1, 0))
	floor.SetMaterial(floor.Material().
		WithTransparency(0.5).
		WithReflective(0.5).
		WithRefractiveIndex(1.5),
	)
	w.AddPrimitives(&floor)

	ball := primitive.NewSphere()
	ball.SetTransform(transform.Translation(0, -3.5, -0.5))
	ball.SetMaterial(ball.Material().
		WithColor(floatcolor.Red).
		WithAmbient(0.5),
	)
	w.AddPrimitives(&ball)

	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	x := primitive.NewIntersection(math.Sqrt2, &floor)

	hc := prepareHitComputations(x, r, x)
	color := w.shadeHit(hc, 5)

	test.AssertAlmost(t, floatcolor.New(0.93391, 0.69643, 0.69243), color)
}

func TestColorWhenRayMisses(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 1, 0))

	c := w.ColorAt(r, 1)

	assert.Equal(t, floatcolor.Black, c)
	assert.EqualValues(t, 1, w.stats.EyeRayCount())
}

func TestColorWhenRayHits(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))

	c := w.ColorAt(r, 1)

	assert.True(t, floatcolor.New(0.38066, 0.47583, 0.2855).Equals(c))
	assert.EqualValues(t, 1, w.stats.EyeRayCount())
}

func TestColorWithIntersectionBehindRay(t *testing.T) {
	w := testWorld()
	outerSphere := w.Primitives()[0]
	outerSphere.SetMaterial(outerSphere.Material().WithAmbient(1))
	innerSphere := w.Primitives()[1]
	innerSphere.SetMaterial(innerSphere.Material().WithAmbient(1))
	r := ray.New(tuple.NewPoint(0, 0, 0.75), tuple.NewVector(0, 0, -1))

	c := w.ColorAt(r, 1)

	assert.True(t, floatcolor.White.Equals(c))
	assert.EqualValues(t, 1, w.stats.EyeRayCount())
}

func TestNoShadowWhenNothingBetweenIntersectionAndLight(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(0, 10, 0)

	assert.False(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestShadowWhenObjectIsBetweenIntersectionAndLight(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(10, -10, 10)

	assert.True(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestNoShadowWhenObjectIsBehindLight(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(-20, 20, -20)

	assert.False(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestNoShadowWhenObjectIsBehindPoint(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(-2, 2, -2)

	assert.False(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestReflectedColorForNonReflectiveMaterial(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	shape := w.primitives[1]
	shape.SetMaterial(shape.Material().WithAmbient(1))
	i := primitive.NewIntersection(1, shape)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, floatcolor.Black, w.reflectedColor(hc, 1))
	assert.EqualValues(t, 0, w.stats.ReflectionRayCount())
}

func TestReflectedColorForReflectiveMaterial(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	shape := primitive.NewPlane()
	shape.SetMaterial(material.Default.WithReflective(0.5))
	shape.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&shape)
	i := primitive.NewIntersection(math.Sqrt2, &shape)

	hc := prepareHitComputations(i, r)

	test.AssertAlmost(t, floatcolor.New(0.19032, 0.2379, 0.14274), w.reflectedColor(hc, 1))
	assert.EqualValues(t, 1, w.stats.ReflectionRayCount())
}

func TestShadeHitForReflectiveMaterial(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	shape := primitive.NewPlane()
	shape.SetMaterial(material.Default.WithReflective(0.5))
	shape.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&shape)
	i := primitive.NewIntersection(math.Sqrt2, &shape)

	hc := prepareHitComputations(i, r)
	color := w.shadeHit(hc, 1)

	test.AssertAlmost(t, floatcolor.New(0.87677, 0.92436, 0.82918), color)
	assert.EqualValues(t, 1, w.stats.ReflectionRayCount())
}

func TestReflectedColorAtMaximumRecursionDepth(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	shape := primitive.NewPlane()
	shape.SetMaterial(material.Default.WithReflective(0.5))
	shape.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&shape)
	i := primitive.NewIntersection(math.Sqrt2, &shape)

	hc := prepareHitComputations(i, r)

	test.AssertAlmost(t, floatcolor.Black, w.reflectedColor(hc, 0))
	assert.EqualValues(t, 0, w.stats.ReflectionRayCount())
}

func TestColorAtWithMutuallyReflectiveSurfaces(t *testing.T) {
	w := New()
	l := light.NewPointLight(tuple.NewPoint(0, 0, 0), floatcolor.White)
	w.AddLights(&l)
	lower := primitive.NewPlane()
	lower.SetMaterial(material.Default.WithReflective(1))
	lower.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&lower)

	upper := primitive.NewPlane()
	upper.SetMaterial(material.Default.WithReflective(1))
	upper.SetTransform(transform.Translation(0, 1, 0))
	w.AddPrimitives(&upper)

	ray := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 1, 0))
	w.ColorAt(ray, 1)

	assert.EqualValues(t, 1, w.stats.ReflectionRayCount())
	// Should not overflow stack from infinite recursion
}

func TestRefractedColorWithOpaqueSurface(t *testing.T) {
	w := testWorld()
	shape := w.primitives[0]
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	xs := primitive.NewIntersections(
		primitive.NewIntersection(4, shape),
		primitive.NewIntersection(6, shape),
	)

	hc := prepareHitComputations(xs[0], r, xs...)
	c := w.refractedColor(hc, 5)

	assert.Equal(t, floatcolor.Black, c)
	assert.EqualValues(t, 0, w.stats.RefractionRayCount())
}

func TestRefractedColorAtMaximumRecursionDepth(t *testing.T) {
	w := testWorld()
	shape := w.primitives[0]
	shape.SetMaterial(shape.Material().
		WithTransparency(1.0).
		WithRefractiveIndex(1.5),
	)
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	xs := primitive.NewIntersections(
		primitive.NewIntersection(4, shape),
		primitive.NewIntersection(6, shape),
	)

	hc := prepareHitComputations(xs[0], r, xs...)
	c := w.refractedColor(hc, 0)

	assert.Equal(t, floatcolor.Black, c)
	assert.EqualValues(t, 0, w.stats.RefractionRayCount())
}

func TestRefractedColorUnderTotalInternalReflection(t *testing.T) {
	w := testWorld()
	shape := w.primitives[0]
	shape.SetMaterial(shape.Material().
		WithTransparency(1.0).
		WithRefractiveIndex(1.5),
	)
	r := ray.New(tuple.NewPoint(0, 0, math.Sqrt2/2), tuple.NewVector(0, 1, 0))
	xs := primitive.NewIntersections(
		primitive.NewIntersection(-math.Sqrt2/2, shape),
		primitive.NewIntersection(math.Sqrt2/2, shape),
	)

	hc := prepareHitComputations(xs[1], r, xs...)
	c := w.refractedColor(hc, 5)

	assert.Equal(t, floatcolor.Black, c)
	assert.EqualValues(t, 0, w.stats.RefractionRayCount())
}

func TestRefractedColorWithRefractedRay(t *testing.T) {
	w := testWorld()
	a := w.primitives[0]
	a.SetMaterial(a.Material().
		WithAmbient(1.0).
		WithPattern(material.TestPattern{}),
	)
	b := w.primitives[1]
	b.SetMaterial(b.Material().
		WithTransparency(1.0).
		WithRefractiveIndex(1.5),
	)
	r := ray.New(tuple.NewPoint(0, 0, 0.1), tuple.NewVector(0, 1, 0))
	xs := primitive.NewIntersections(
		primitive.NewIntersection(-0.9899, a),
		primitive.NewIntersection(-0.4899, b),
		primitive.NewIntersection(0.4899, b),
		primitive.NewIntersection(0.9899, a),
	)

	hc := prepareHitComputations(xs[2], r, xs...)
	c := w.refractedColor(hc, 5)

	test.AssertAlmost(t, floatcolor.New(0, 0.99888, 0.04725), c)
	assert.EqualValues(t, 1, w.stats.RefractionRayCount())
}

func TestSchlickApproximationUnderTotalInternalReflection(t *testing.T) {
	s := glassSphere()
	r := ray.New(tuple.NewPoint(0, 0, math.Sqrt2/2), tuple.NewVector(0, 1, 0))
	xs := primitive.NewIntersections(primitive.NewIntersection(-math.Sqrt2/2, s), primitive.NewIntersection(math.Sqrt2/2, s))

	hc := prepareHitComputations(xs[1], r, xs...)

	reflectance := schlick(hc)

	assert.Equal(t, 1.0, reflectance)
}

func TestSchlickApproximationWithPerpendicularViewingAngle(t *testing.T) {
	s := glassSphere()
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 1, 0))
	xs := primitive.NewIntersections(primitive.NewIntersection(-1, s), primitive.NewIntersection(1, s))

	hc := prepareHitComputations(xs[1], r, xs...)

	reflectance := schlick(hc)

	test.AssertAlmost(t, 0.04, reflectance)
}

func TestSchlickApproximationWithSmallAngle(t *testing.T) {
	s := glassSphere()
	r := ray.New(tuple.NewPoint(0, 0.99, -2), tuple.NewVector(0, 0, 1))
	xs := primitive.NewIntersections(primitive.NewIntersection(1.8589, s))

	hc := prepareHitComputations(xs[0], r, xs...)

	reflectance := schlick(hc)

	test.AssertAlmost(t, 0.48873, reflectance)
}

func TestGenerateID(t *testing.T) {
	w := New()
	numIDs := 10000
	ch := make(chan ID, numIDs)
	var wg sync.WaitGroup
	wg.Add(numIDs)

	for i := 0; i < numIDs; i++ {
		go func() {
			ch <- w.NextID()
			wg.Done()
		}()
	}

	wg.Wait()

	last := <-ch
	for i := 1; i < numIDs; i++ {
		next := <-ch
		t.Logf("checking %v and %v\n", last, next)
		require.NotEqual(t, last, next)
		last = next
	}
}

func testWorld() *World {
	w := New()

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)
	s1 := primitive.NewSphere()
	s1.SetMaterial(material.Default.
		WithColor(floatcolor.New(0.8, 1.0, 0.6)).
		WithDiffuse(0.7).
		WithSpecular(0.2),
	)

	s2 := primitive.NewSphere()
	s2.SetTransform(transform.Scaling(0.5, 0.5, 0.5))

	w.AddLights(&light)
	w.AddPrimitives(&s1, &s2)

	return w
}

func glassSphere() primitive.Primitive {
	s := primitive.NewSphere()
	s.SetMaterial(material.Default.
		WithTransparency(1).
		WithRefractiveIndex(1.5),
	)
	return &s
}
