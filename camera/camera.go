package camera

import (
	"image"
	"math"
	"runtime"
	"sync"

	"github.com/danieltmartin/ray-tracer/canvas"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/danieltmartin/ray-tracer/world"
)

const recursionDepth = 5

type Camera struct {
	hsize            uint
	vsize            uint
	fieldOfView      float64
	transform        matrix.Matrix
	inverseTransform matrix.Matrix
	pixelSize        float64
	halfWidth        float64
	halfHeight       float64
}

func New(hsize, vsize uint, fieldOfView float64) *Camera {
	c := &Camera{
		hsize:            hsize,
		vsize:            vsize,
		fieldOfView:      fieldOfView,
		transform:        matrix.Identity4(),
		inverseTransform: matrix.Identity4(),
	}
	computePixelSizeAndDimensions(c)
	return c
}

func (c *Camera) SetTransform(t matrix.Matrix) {
	c.transform = t
	c.inverseTransform = t.Inverse()
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
	cameraTransform := c.inverseTransform
	pixel := cameraTransform.MulTuple(untransformedPixel)
	origin := cameraTransform.MulTuple(tuple.NewPoint(0, 0, 0))
	direction := pixel.Sub(origin).Norm()

	return ray.New(origin, direction)
}

func (c *Camera) Render(w *world.World) image.Image {
	canvas := canvas.New(c.hsize, c.vsize)

	semaphore := make(chan bool, runtime.NumCPU())
	var wg sync.WaitGroup

	for y := uint(0); y < c.vsize; y++ {
		semaphore <- true // Acquire
		wg.Add(1)

		go func(y uint) {
			defer wg.Done()
			for x := uint(0); x < c.hsize; x++ {
				r := c.RayForPixel(x, y)
				color := w.ColorAt(r, recursionDepth)
				canvas.WritePixel(x, y, color)
			}
			<-semaphore // Release
		}(y)
	}

	wg.Wait()

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
