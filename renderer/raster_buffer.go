package renderer

// RasterBuffer provides a memory mapped RGBA and Z buffer
// This buffer must be blitted to another buffer, for example,
// PNG or display buffer (like SDL).
type RasterBuffer struct {
	colorBuf [][]RGBA
	zBuf     [][]float32

	width      int
	height     int
	ClearDepth float32
	ClearColor RGBA

	PixelColor RGBA
}

// NewRasterBuffer creates a display buffer
func NewRasterBuffer(width, height int) *RasterBuffer {
	db := new(RasterBuffer)
	db.width = width
	db.height = height

	db.ClearDepth = -100000000.0 // Default value
	db.ClearColor.R = 127
	db.ClearColor.G = 127
	db.ClearColor.B = 127
	db.ClearColor.A = 255

	db.colorBuf = make([][]RGBA, width)
	for i := range db.colorBuf {
		db.colorBuf[i] = make([]RGBA, height)
	}
	db.zBuf = make([][]float32, width)
	for i := range db.zBuf {
		db.zBuf[i] = make([]float32, height)
	}
	return db
}

// Clear clears only the color buffer
func (rb *RasterBuffer) Clear() {
	for i := range rb.colorBuf {
		for j := 0; j < rb.height; j++ {
			rb.colorBuf[i][j].R = rb.ClearColor.R
			rb.colorBuf[i][j].G = rb.ClearColor.G
			rb.colorBuf[i][j].B = rb.ClearColor.B
			rb.colorBuf[i][j].A = rb.ClearColor.A
		}
	}
}

// ClearZs sets the z buffer to ClearDepth
func (rb *RasterBuffer) ClearZs() {
	for i := range rb.colorBuf {
		for j := 0; j < rb.height; j++ {
			rb.zBuf[i][j] = rb.ClearDepth
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
		rb.colorBuf[x][y] = rb.PixelColor
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
func (rb *RasterBuffer) SetPixelColor(x, y int, z float32, color RGBA) int {
	rb.PixelColor = color
	return rb.SetPixel(x, y, z)
}

func (rb *RasterBuffer) drawLine(xP, yP, xQ, yQ int, zP, zQ float32, color RGBA) {
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
			if x > rb.width || y > rb.height {
				break // pixel off screen
			}
			zstatus := rb.SetPixel(x, y, z)
			if zstatus == 1 {
				rb.colorBuf[x][y] = rb.PixelColor
			}
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
			if x > rb.width || y > rb.height {
				break // pixel off screen
			}
			zstatus := rb.SetPixel(x, y, z)
			if zstatus == 1 {
				rb.colorBuf[x][y] = rb.PixelColor
			}
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
