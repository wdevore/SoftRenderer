package graphics

// Polygon is specified using triangles where the vertices
// are shared with other triangles that represent the polygon.
// During each
type Polygon struct {
	vertices []int // Indices into vertex buffer
}
