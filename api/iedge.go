package api

// IEdge is an edge for drawing a horizontal line.
type IEdge interface {
	Set(x1, y1, x2, y2 int, z1, z2 float32)
	Step() bool

	XY() (x, y int)
	YBot() int
	Z1() float32
	Z2() float32
}
