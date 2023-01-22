package camera

import (
	"image"
	"math"

	"github.com/danieltmartin/ray-tracer/canvas"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/danieltmartin/ray-tracer/world"
)

type Camera struct {
	hsize       uint
	vsize       uint
	fieldOfView float64
	transform   matrix.Matrix
	pixelSize   float64
	halfWidth   float64
	halfHeight  float64
}

func New(hsize, vsize uint, fieldOfView float64) *Camera {
	c := &Camera{
		hsize:       hsize,
		vsize:       vsize,
		fieldOfView: fieldOfView,
		transform:   matrix.Identity4(),
	}
	computePixelSizeAndDimensions(c)
	return c
}

func (c *Camera) SetTransform(t matrix.Matrix) {
	c.transform = t
}

func (c *Camera) RayForPixel(px, py uint) ray.Ray {
	// Offset from edge of canvas to center of pixel
	xOffset := (float64(px) + 0.5) * c.pixelSize
	yOffset := (float64(py) + 0.5) * c.pixelSize

	untransformedPixel := tuple.NewPoint(
		c.halfWidth-xOffset,
		c.halfHeight-yOffset,
		-1,
	)

	// Transform the pixel to account for camera position.
	// We use the inverse because the transform matrix transforms the world, not the camera,
	// so the inverse is the transform of the camera.
	cameraTransform := c.transform.Inverse()
	pixel := cameraTransform.MulTuple(untransformedPixel)
	origin := cameraTransform.MulTuple(tuple.NewPoint(0, 0, 0))
	direction := pixel.Sub(origin).Norm()

	return ray.New(origin, direction)
}

func (c *Camera) Render(w *world.World) image.Image {
	canvas := canvas.New(c.hsize, c.vsize)

	for y := uint(0); y < c.vsize; y++ {
		for x := uint(0); x < c.hsize; x++ {
			r := c.RayForPixel(x, y)
			color := w.ColorAt(r)
			canvas.WritePixel(x, y, color)
		}
	}

	return canvas
}

func computePixelSizeAndDimensions(c *Camera) {
	halfView := math.Tan(c.fieldOfView / 2)
	aspectRatio := float64(c.hsize) / float64(c.vsize)
	if aspectRatio >= 1 {
		c.halfWidth = halfView
		c.halfHeight = halfView / aspectRatio
	} else {
		c.halfHeight = halfView
		c.halfWidth = halfView * aspectRatio
	}

	c.pixelSize = c.halfWidth * 2 / float64(c.hsize)
}
