package graphics

import "SoftRenderer/api"

// Edge is part of a triangle
type Edge struct {
	// The triangle this edge is associated with
	triangle api.ITriangle

	// Vertex indices into the ITriangle
	p, q int

	// Is this edge shared with another edge. If it's then
	// it is only rendered if it's a left or top edge--unless
	// the visibility flag is set.
	shared bool

	// Forces rendering of the edge regardless if it's shared.
	visible bool
}
