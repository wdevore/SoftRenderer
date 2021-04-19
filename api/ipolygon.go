package api

// IPolygon represents a collection of triangles that share edges.
// The edges belong to triangles. The outer edges are never shared but all
// internal ones are.
// When triangle are split internally the top/bottom are shared edges
// of which one is drawn.
// In the top-left rule only the top/left edges are drawn which means the
// the bottom/right edges aren't--with the exception being if the edge is
// an outer edge. Outer edge pixels are always drawn.

type IPolygon interface {
	AddVertex(x, y int, z float32)
	AddTriangle(x1, y1, x2, y2, x3, y3 int, sharedE1, sharedE2, sharedE3 bool)
	Draw(raster IRasterBuffer)
	Fill(raster IRasterBuffer)
}
