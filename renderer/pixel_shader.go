package renderer

import "SoftRenderer/api"

// PixelShader is a psuedo shader meaning that it isn't programmable. It just
// represents the pixel rasterization part of a pipeline.
// This shader works on left/right edges. Each pair represents a triangle with
// either a flat-top or flat-bottom.
type PixelShader struct {
	// Edge buffer is a pre-allocated buffer of fixed size.
	// As a triangle is processed two edges are produced and placed into
	// this buffer
	edges []api.IEdge // Pairs of edges
}

func NewPixeShader() api.IPixelShader {
	o := new(PixelShader)
	return o
}
