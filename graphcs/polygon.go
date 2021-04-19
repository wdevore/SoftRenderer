package graphics

import "SoftRenderer/api"

// Polygon is specified using triangles where the vertices
// are shared with other triangles that represent the polygon.

// Our polygon can only be defined as Convex meaning there can't
// be any insets or cavities.
type Polygon struct {
	vertices []int // Indices into vertex buffer
}

func NewPolygon() api.IPolygon {
	o := new(Polygon)
	return o
}

func (p *Polygon) AddVertex(x, y int, z float32) {

}

func (p *Polygon) AddTriangle(x1, y1, x2, y2, x3, y3 int, sharedE1, sharedE2, sharedE3 bool) {

}

func (p *Polygon) Draw(raster api.IRasterBuffer) {

}

func (p *Polygon) Fill(raster api.IRasterBuffer) {

}
