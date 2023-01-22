package canvas

import (
	"image"
	"image/color"

	"github.com/danieltmartin/ray-tracer/floatcolor"
)

type Canvas struct {
	pixels []floatcolor.Float64Color
	width  uint
	height uint
}

func New(width, height uint) Canvas {
	pixels := make([]floatcolor.Float64Color, width*height)
	return Canvas{
		pixels, width, height,
	}
}

func (c Canvas) Width() uint {
	return c.width
}

func (c Canvas) Height() uint {
	return c.height
}

func (c Canvas) PixelAt(x, y uint) floatcolor.Float64Color {
	return c.pixels[x+y*c.width]
}

func (c Canvas) WritePixel(x, y uint, col floatcolor.Float64Color) {
	if x > c.width-1 || y > c.height-1 {
		return
	}
	c.pixels[x+y*c.width] = col
}

func (c Canvas) ColorModel() color.Model {
	return floatcolor.Float64Model
}

func (c Canvas) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(c.width), int(c.height))
}

func (c Canvas) At(x, y int) color.Color {
	return c.PixelAt(uint(x), uint(y))
}
