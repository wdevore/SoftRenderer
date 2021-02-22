package api

// IRasterizer api for line and triangle rasterization
type IRasterizer interface {
	DrawLine(surface ISurface, x1, y1, x2, y2 int)

	DrawLineAmmeraal(surface ISurface, direction bool, x1, y1, x2, y2 int)

	DrawLineAmmeraalDxGtDy(surface ISurface, up bool, x1, y1, x2, y2 int)
	DrawLineAmmeraalDyGtDx(surface ISurface, up bool, x1, y1, x2, y2 int)

	DrawLineAmmeraalDxGtDy2(surface ISurface, up bool, x1, y1, x2, y2 int)
	DrawLineAmmeraalDyGtDx2(surface ISurface, up bool, x1, y1, x2, y2 int)

	DrawTriangle(surface ISurface, x1, y1, x2, y2, x3, y3 int)
	RenderTriangle(surface ISurface, x1, y1, x2, y2, x3, y3 int)
	FillTriangle(surface ISurface, flatTop bool, x1, y1, x2, y2, x3, y3 int)

	Sort(x1, y1, x2, y2, x3, y3 *int)
}
