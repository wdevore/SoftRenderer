package renderer

import (
	"SoftRenderer/api"
	"image"
	"image/color"
)

// RasterBuffer provides a memory mapped RGBA and Z buffer
// This buffer must be blitted to another buffer, for example,
// PNG or display buffer (like SDL).
type RasterBuffer struct {
	width  int
	height int

	// Image pixels
	pixels *image.RGBA
	bounds image.Rectangle

	// ZBuffer
	ClearDepth    float32
	zBuf          [][]float32
	alphaBlending bool

	// Pen colors
	ClearColor color.RGBA
	PixelColor color.RGBA
}

// NewRasterBuffer creates a display buffer
func NewRasterBuffer(width, height int) api.IRasterBuffer {
	o := new(RasterBuffer)
	o.width = width
	o.height = height

	o.alphaBlending = false

	o.bounds = image.Rect(0, 0, width, height)
	o.pixels = image.NewRGBA(o.bounds)

	o.ClearColor.R = 127
	o.ClearColor.G = 127
	o.ClearColor.B = 127
	o.ClearColor.A = 255

	o.ClearDepth = -100000000.0 // Default value
	o.zBuf = make([][]float32, width)
	for i := range o.zBuf {
		o.zBuf[i] = make([]float32, height)
	}

	return o
}

// EnableAlphaBlending turns on/off per pixel alpha blending
func (rb *RasterBuffer) EnableAlphaBlending(enable bool) {
	rb.alphaBlending = enable
}

// Pixels returns underlying color buffer
func (rb *RasterBuffer) Pixels() *image.RGBA {
	return rb.pixels
}

// Clear clears both color and depth buffers
func (rb *RasterBuffer) Clear() {
	for y := 0; y < rb.height; y++ {
		for x := 0; x < rb.width; x++ {
			rb.pixels.SetRGBA(x, y, rb.ClearColor)
			rb.zBuf[x][y] = rb.ClearDepth
		}
	}
}

// ClearColorBuffer clears only the color/pixel buffer
func (rb *RasterBuffer) ClearColorBuffer() {
	/// TODO use image/draw to clear using a SRC
	for y := 0; y < rb.height; y++ {
		for x := 0; x < rb.width; x++ {
			rb.pixels.SetRGBA(x, y, rb.ClearColor)
		}
	}
}

// ClearDepthBuffer sets the z buffer to ClearDepth
func (rb *RasterBuffer) ClearDepthBuffer() {
	for y := 0; y < rb.height; y++ {
		for x := 0; x < rb.width; x++ {
			rb.zBuf[x][y] = rb.ClearDepth
		}
	}
}

// SetPixel sets a pixel and rejects based on a Z buffer.
// Returns:
// -1 = pixel is beyond screen
// 0 = pixel was farther away and ignored
// 1 = pixel is closer and was entered into framebuffer and zbuffer
// 2 = pixel is exact/(on top) and was ignored
func (rb *RasterBuffer) SetPixel(x, y int, z float32) int {
	if x < 0 || x > rb.width || y < 0 || y > rb.height {
		return -1
	}

	zd := rb.zBuf[x][y]

	if z < zd {
		//////////////////////////////////
		// pixel farther away
		//////////////////////////////////
		return 0
	} else if z > zd {
		//////////////////////////////////
		// pixel closer (i.e. on top and visible)
		//////////////////////////////////
		rb.zBuf[x][y] = z

		// https://en.wikipedia.org/wiki/Alpha_compositing Alpha blending section
		// Non premultiplied alpha
		if rb.alphaBlending {
			dst := rb.pixels.RGBAAt(x, y)
			src := rb.PixelColor
			A := float32(src.A) / 255.0
			dst.R = uint8(float32(src.R)*A + float32(dst.R)*(1.0-A))
			dst.G = uint8(float32(src.G)*A + float32(dst.G)*(1.0-A))
			dst.B = uint8(float32(src.B)*A + float32(dst.B)*(1.0-A))
			dst.A = 255
			rb.pixels.SetRGBA(x, y, dst)
		} else {
			rb.pixels.SetRGBA(x, y, rb.PixelColor)
		}

		return 1
	} else {
		//////////////////////////////////
		// pixel same distance
		//////////////////////////////////
		return 2
	}
}

// SetPixelColor set the current pixel color and sets the pixel
// using SetPixel()
func (rb *RasterBuffer) SetPixelColor(c color.RGBA) {
	rb.PixelColor = c
}

// DrawLine draws a line into the buffer
func (rb *RasterBuffer) DrawLine(xP, yP, xQ, yQ int, zP, zQ float32) {
	if xP < 0 || xP > rb.width-1 || xQ < 0 || xQ > rb.width-1 {
		return
	}
	if yP < 0 || yP > rb.height-1 || yQ < 0 || yQ > rb.height-1 {
		return
	}

	zrP := 1.0 / zP
	zrQ := 1.0 / zQ
	dzr := zrQ - zrP
	dx := (xQ - xP)
	dy := (yQ - yP)
	z := zrP
	dzdx := dzr / float32(dx)
	dzdy := dzr / float32(dy)

	x := xP
	y := yP
	D := 0
	HX := xQ - xP
	HY := yQ - yP
	xInc := 1
	yInc := 1

	if HX < 0 {
		xInc = -1
		dzdx = -dzdx
		HX = -HX
	}
	if HY < 0 {
		yInc = -1
		dzdy = -dzdy
		HY = -HY
	}
	if HY <= HX {
		c := 2 * HX
		M := 2 * HY
		for {
			rb.SetPixel(x, y, z)
			if x == xQ {
				break
			}
			x += xInc
			z += dzdx
			D += M
			if D > HX {
				y += yInc
				D -= c
			}
		}
	} else {
		c := 2 * HY
		M := 2 * HX
		for {
			rb.SetPixel(x, y, z)
			if y == yQ {
				break
			}
			y += yInc
			z += dzdy
			D += M
			if D > HY {
				x += xInc
				D -= c
			}
		}
	}
}

// DrawLineAmmeraal has no zbuffer support
func (rb *RasterBuffer) DrawLineAmmeraal(xP, yP, xQ, yQ int, zP, zQ float32) {
	x := xP
	y := yP
	d := 0

	yInc := 1
	xInc := 1
	dx := xQ - xP
	dy := yQ - yP

	if dx < 0 {
		xInc = -1
		dx = -dx
	}
	if dy < 0 {
		yInc = -1
		dy = -dy
	}

	// --------------------------------------------------------------------
	if dy <= dx {
		// The change is Y is smaller than X
		m := dy << 1
		c := dx << 1

		if xInc < 0 {
			dx++
		}

		col := uint8(0)
		for true {
			// rb.SetPixelColor(color.RGBA{R: 0, G: col, B: col, A: 255})
			rb.SetPixel(x, y, zP)
			col += 3

			if x == xQ {
				break
			}

			// X is the major step axis
			x += xInc
			d += m
			if d >= dx {
				y += yInc
				d -= c
			}
		}
	} else {
		c := dy << 1
		m := dx << 1

		if yInc < 0 {
			dy++
		}

		col := uint8(0)
		for true {
			// rb.SetPixelColor(color.RGBA{R: col, G: 0, B: 0, A: 255})
			rb.SetPixel(x, y, zP)
			col += 3

			if y == yQ {
				break
			}

			// Y is the major step axis
			y += yInc
			d += m
			if d >= dy {
				x += xInc
				d -= c
			}
		}
	}
}

// FillTriangleAmmeraal --
func (rb *RasterBuffer) FillTriangleAmmeraal(leftEdge, rightEdge api.IEdge, skipLast bool) {
	lx, ly := leftEdge.XY()
	rx, ry := rightEdge.XY()

	dy := 0
	if lx > rx {
		dy = -1
	}

	if skipLast && ly == leftEdge.YBot()-dy {
		return
	}

	for x := lx; x <= rx; x++ {
		rb.SetPixel(x, ly, leftEdge.Z1())
	}

	for leftEdge.Step() {
		lx, ly = leftEdge.XY()

		// For slopes where dy > dx, ry needs to "catch up" to ly
		// because ly is changing on each step where ry isn't.
		for ry < ly {
			rightEdge.Step()
			rx, ry = rightEdge.XY()
		}

		// If the "side" vertex is to the right then there isn't
		// a line to skip.
		dy = 0
		if lx > rx {
			dy = -1
		}

		if skipLast && ly == leftEdge.YBot()-dy {
			return
		}

		// We always want to fill the scanline from left to right
		if lx > rx {
			t := rx
			rx = lx
			lx = t
		}

		// Fill scanline
		for x := lx; x <= rx; x++ {
			rb.SetPixel(x, ly, leftEdge.Z1())
		}
	}
}
