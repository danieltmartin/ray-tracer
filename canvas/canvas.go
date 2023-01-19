package canvas

import (
	"image"
	"image/color"

	"github.com/danieltmartin/ray-tracer/floatcolor"
)

type Canvas struct {
	pixels []floatcolor.Float64Color
	width  int
	height int
}

func New(width, height int) Canvas {
	pixels := make([]floatcolor.Float64Color, width*height)
	return Canvas{
		pixels, width, height,
	}
}

func (c Canvas) Width() int {
	return c.width
}

func (c Canvas) Height() int {
	return c.height
}

func (c Canvas) PixelAt(x, y int) floatcolor.Float64Color {
	return c.pixels[x+y*c.width]
}

func (c Canvas) WritePixel(x, y int, col floatcolor.Float64Color) {
	c.pixels[x+y*c.width] = col
}

func (c Canvas) ColorModel() color.Model {
	return floatcolor.Float64Model
}

func (c Canvas) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.width, c.height)
}

func (c Canvas) At(x, y int) color.Color {
	return c.PixelAt(x, y)
}
