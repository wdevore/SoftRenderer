package renderer

import "SoftRenderer/api"

// VertexShader is a psuedo shader meaning that it isn't programmable. It just
// represents the geometry part of a pipeline
type VertexShader struct {
	// Vertex buffer holds all vertices.
	// They are transformed into another vertex pipeline buffer later in
	// the pipeline.
	// Format is: x,y,z,x,y,z,x,y,z...
	vertices []float32

	index int
}

func NewVertexShader() api.IVertexShader {
	o := new(VertexShader)
	o.index = 0
	return o
}

func (vs *VertexShader) AddVertex(x, y, z float32) int {
	vs.vertices = append(vs.vertices, x, y, z)
	i := vs.index
	vs.index++
	return i
}

func (vs *VertexShader) Buffer() []float32 {
	return vs.vertices
}

// Transform
