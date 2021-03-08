package api

import (
	"image"
	"image/color"
)

// IRasterBuffer api for color and depth buffer
type IRasterBuffer interface {
	EnableAlphaBlending(enable bool)
	Pixels() *image.RGBA
	Clear()
	SetPixel(x, y int, z float32) int
	SetPixelColor(c color.RGBA)

	DrawLine(xP, yP, xQ, yQ int, zP, zQ float32)
	DrawLineAmmeraal(xP, yP, xQ, yQ int, zP, zQ float32)

	FillTriangleAmmeraal(leftEdge, rightEdge IEdge, skipLast bool)
}
