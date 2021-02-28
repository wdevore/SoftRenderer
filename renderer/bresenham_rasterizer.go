package renderer

import (
	"SoftRenderer/api"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

// SDL2 coordinate space
//
//     -Y
//     ^
//     |
//     |
//     |
// 0,0 .-------> +X
//     |
//     |
//     |
//     v
//     +Y

type bresenhamLineRasterizer struct {
	x    int
	y    int
	c    int
	M    int
	D    int
	DX   int
	DY   int
	xInc int
	yInc int
}

// NewBresenHamRasterizer --
func NewBresenHamRasterizer() api.IRasterizer {
	o := new(bresenhamLineRasterizer)
	return o
}

func (r *bresenhamLineRasterizer) DrawLineAmmeraal(buffer api.IRasterBuffer, direction bool, xP, yP, xQ, yQ int) {
	// P -> Q or Q -> P
	if direction {
		tx := xP
		ty := yP
		xP = xQ
		yP = yQ
		xQ = tx
		yQ = ty
	}

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
		m := dy << 1
		c := dx << 1

		if xInc < 0 {
			dx++
		}

		col := uint8(0)
		for true {
			buffer.SetPixelColor(color.RGBA{R: 0, G: 0, B: col, A: 255})
			buffer.SetPixel(x, y, 0.0)
			col += 2

			if x == xQ {
				break
			}

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
			buffer.SetPixelColor(color.RGBA{R: col, G: 0, B: 0, A: 255})
			buffer.SetPixel(x, y, 0.0)
			col += 2

			if y == yQ {
				break
			}

			y += yInc
			d += m
			if d >= dy {
				x += xInc
				d -= c
			}
		}
	}
}

// For Horizontal and between -45 to 45 and 135 to 225
//       \    /
//        \  /
//   Cyan  \/  Cyan
//         /\
//        /  \
//       /    \
func (r *bresenhamLineRasterizer) DrawLineAmmeraalDyGtDx(surface api.ISurface, down bool, xP, yP, xQ, yQ int) {
	if !down {
		tx := xP
		ty := yP
		xP = xQ
		yP = yQ
		xQ = tx
		yQ = ty
	}

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

	// ------------------------------------------------------
	c := dy << 1
	m := dx << 1

	if yInc < 0 {
		dy++
	}

	col := uint8(0)
	for true {
		surface.SetDrawColor(sdl.Color{R: col, G: col, B: col, A: 255})
		surface.SetPixel(x, y)
		col++

		if y == yQ {
			break
		}

		y += yInc //
		d += m
		if d >= dy {
			x += xInc
			d -= c
		}
	}
}

func (r *bresenhamLineRasterizer) DrawLineAmmeraalDxGtDy2(surface api.ISurface, down bool, xP, yP, xQ, yQ int) {
	if down {
		tx := xP
		ty := yP
		xP = xQ
		yP = yQ
		xQ = tx
		yQ = ty
	}

	x := xP
	y := yP
	d := 0
	dx := 0
	m := 0
	inc := 1

	if down {
		dx = -(xQ - xP)       // <--
		m = (-(yQ - yP)) << 1 // <--
		inc = -1
	} else {
		dx = (xQ - xP)     // <--
		m = (yQ - yP) << 1 // <--
	}

	c := dx << 1

	col := uint8(0)
	for true {
		surface.SetDrawColor(sdl.Color{R: col, G: col, B: col, A: 255})
		surface.SetPixel(x, y)
		col++

		if x == xQ {
			break
		}
		x += inc // <--
		d += m
		if d >= dx {
			y += inc // <--
			d -= c
		}
	}
}

func (r *bresenhamLineRasterizer) DrawLineAmmeraalDyGtDx2(surface api.ISurface, down bool, xP, yP, xQ, yQ int) {
	if !down {
		tx := xP
		ty := yP
		xP = xQ
		yP = yQ
		xQ = tx
		yQ = ty
	}

	x := xP
	y := yP
	d := 0

	yInc := 1
	xInc := -1
	dx := -(xQ - xP)
	dy := (yQ - yP)

	if down {
		yInc = -yInc
		xInc = -xInc
		dx = -dx
		dy = -dy
	}

	c := dy << 1
	m := dx << 1

	if yInc < 0 {
		dy++
	}

	col := uint8(0)
	for true {
		surface.SetDrawColor(sdl.Color{R: col, G: col, B: col, A: 255})
		surface.SetPixel(x, y)
		col++

		if y == yQ {
			break
		}

		y += yInc
		d += m
		if d >= dy {
			x += xInc
			d -= c
		}
	}

}

// For Horizontal and between 45 to 135 and 225 to 270
//    135       45
//      \Yellow/
//       \    /
//        \  /
//         \/     ____ 0 degrees
//         /\
//        /  \
//       /    \
//      /Yellow\
func (r *bresenhamLineRasterizer) DrawLineAmmeraalDxGtDy(surface api.ISurface, down bool, xP, yP, xQ, yQ int) {
	if !down {
		tx := xP
		ty := yP
		xP = xQ
		yP = yQ
		xQ = tx
		yQ = ty
	}

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

	// ------------------------------------------------------
	m := dy << 1
	c := dx << 1

	if xInc < 0 {
		dx++
	}

	col := uint8(0)
	for true {
		surface.SetDrawColor(sdl.Color{R: col, G: col, B: col, A: 255})
		surface.SetPixel(x, y)
		col++

		if x == xQ {
			break
		}

		x += xInc
		d += m
		if d >= dx {
			y += yInc
			d -= c
		}
	}
}

func (r *bresenhamLineRasterizer) DrawLine(surface api.ISurface, x1, y1, x2, y2 int) {
	r.x = x1
	r.y = y1
	r.D = 0
	r.DX = x2 - x1
	r.DY = y2 - y1
	r.xInc = 1
	r.yInc = 1

	if r.DX < 0 {
		r.xInc = -1
		r.DX = -r.DX
	}

	if r.DY < 0 {
		r.yInc = -1
		r.DY = -r.DY
	}

	if r.DY <= r.DX {
		r.c = r.DX << 1 //2 * HX
		r.M = r.DY << 1 //2 * HY
		for true {
			// For Horizontal and between -45 to 45 and 135 to 225
			//       \    /
			//        \  /
			//   Cyan  \/  Cyan
			//         /\
			//        /  \
			//       /    \
			// surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 255, A: 255}) // Cyan
			surface.SetPixel(r.x, r.y)

			if r.x == x2 {
				break
			}

			r.x += r.xInc
			r.D += r.M
			if r.D > r.DX {
				r.y += r.yInc
				r.D -= r.c
			}
		}
	} else {
		r.c = r.DY << 1 //2 * HY
		r.M = r.DX << 1 //2 * HX
		for true {
			// For vertical and between 45 to 135 and 225 to 270
			//      \Yellow/
			//       \    /
			//        \  /
			//         \/
			//         /\
			//        /  \
			//       /    \
			//      /Yellow\
			// surface.SetDrawColor(sdl.Color{R: 255, G: 255, B: 0, A: 255}) // Yellow
			surface.SetPixel(r.x, r.y)

			if r.y == y2 {
				break
			}

			r.y += r.yInc
			r.D += r.M
			if r.D > r.DY {
				r.x += r.xInc
				r.D -= r.c
			}
		}
	}
}

// Draws an outline only
func (r *bresenhamLineRasterizer) DrawTriangle(surface api.ISurface, x1, y1, x2, y2, x3, y3 int) {
	// sort the three vertices by y-coordinate ascending,
	// so x1,y1 is the topmost (max y) vertex
	r.Sort(&x1, &y1, &x2, &y2, &x3, &y3)

	if y2 == y3 {
		// Case for flat-bottom triangle
		surface.SetDrawColor(sdl.Color{R: 255, G: 127, B: 0, A: 255})
		r.DrawLine(surface, x1, y1, x2, y2) // Diagonal
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x3, y3) // Bottom
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x3, y3, x1, y1) // Left
	} else if y1 == y2 {
		// Case for flat-top triangle
		surface.SetDrawColor(sdl.Color{R: 255, G: 127, B: 0, A: 255})
		r.DrawLine(surface, x3, y3, x1, y1) // Diagonal
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x1, y1, x2, y2) // Top
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x3, y3) // Left
	} else {
		// General case
		// split the triangle into two triangles: top-half and bottom-half
		x := int(float32(x1) + (float32(y2-y1)/float32(y3-y1))*float32(x3-x1))

		// flat-bottom
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x1, y1, x2, y2) // Left
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x, y2) // Bottom
		surface.SetDrawColor(sdl.Color{R: 0, G: 0, B: 255, A: 255})
		r.DrawLine(surface, x, y2, x1, y1) // Right

		// flat-top
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x3, y3, x2, y2) // Left
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x, y2) // Top
		surface.SetDrawColor(sdl.Color{R: 0, G: 0, B: 255, A: 255})
		r.DrawLine(surface, x, y2, x3, y3) // Right
	}
}

// Rasterizes it via FillTriangle
func (r *bresenhamLineRasterizer) RenderTriangle(surface api.ISurface, x1, y1, x2, y2, x3, y3 int) {
	// sort the three vertices by y-coordinate ascending,
	// so x1,y1 is the topmost (max y) vertex
	r.Sort(&x1, &y1, &x2, &y2, &x3, &y3)

	// Now y1 <= y2 <= y3
	// With this knowledge we can now rasterize using edges

	if y2 == y3 {
		// Case for flat-bottom triangle
		surface.SetDrawColor(sdl.Color{R: 255, G: 127, B: 0, A: 255})
		r.DrawLine(surface, x1, y1, x2, y2) // Diagonal
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x3, y3) // Bottom
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x3, y3, x1, y1) // Left
	} else if y1 == y2 {
		// Case for flat-top triangle
		surface.SetDrawColor(sdl.Color{R: 255, G: 127, B: 0, A: 255})
		r.DrawLine(surface, x3, y3, x1, y1) // Diagonal
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x1, y1, x2, y2) // Top
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x3, y3) // Left
	} else {
		// General case
		// split the triangle into two triangles: top-half and bottom-half
		x := int(float32(x1) + (float32(y2-y1)/float32(y3-y1))*float32(x3-x1))

		// flat-bottom
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x1, y1, x2, y2) // Left
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x, y2) // Bottom
		surface.SetDrawColor(sdl.Color{R: 0, G: 0, B: 255, A: 255})
		r.DrawLine(surface, x, y2, x1, y1) // Right

		// flat-top
		surface.SetDrawColor(sdl.Color{R: 255, G: 0, B: 0, A: 255})
		r.DrawLine(surface, x3, y3, x2, y2) // Left
		surface.SetDrawColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
		r.DrawLine(surface, x2, y2, x, y2) // Top
		surface.SetDrawColor(sdl.Color{R: 0, G: 0, B: 255, A: 255})
		r.DrawLine(surface, x, y2, x3, y3) // Right
	}
}

func (r *bresenhamLineRasterizer) FillTriangle(surface api.ISurface, flatTop bool, x1, y1, x2, y2, x3, y3 int) {
	// This fill algorithm always horizontal lines between edges
	// The right most pixel is NOT drawn (aka the right edge)

	// Is this triangle a flat-top or flat-bottom
	if flatTop {
		// flat-top: This triangle has a left AND top edge that is drawn.
		// It only has a right edge that is NOT drawn.
	} else {
		// flat-bottom: This triangle only has a left edge that is drawn.
		// The bottom edge is NOT drawn.
	}
}

// Sorts as y1 <= y2 <= y3
func (r *bresenhamLineRasterizer) Sort(x1, y1, x2, y2, x3, y3 *int) {
	x := 0
	y := 0
	// fmt.Println("Before: ", *x1, *y1, *x2, *y2, *x3, *y3)

	// Make y1 <= y2 if needed
	if *y1 > *y2 {
		x = *x1
		y = *y1
		*x1 = *x2
		*y1 = *y2
		*x2 = x
		*y2 = y
	}

	// Now y1 <= y2. Make y1 <= y3
	if *y1 > *y3 {
		x = *x1
		y = *y1
		*x1 = *x3
		*y1 = *y3
		*x3 = x
		*y3 = y
	}

	// Now y1 <= y2 and y1 <= y3. Make y2 <= y3
	if *y2 > *y3 {
		x = *x2
		y = *y2
		*x2 = *x3
		*y2 = *y3
		*x3 = x
		*y3 = y
	}

	// fmt.Println("After: ", *x1, *y1, *x2, *y2, *x3, *y3)
}
