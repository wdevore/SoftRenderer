package api

// ITriangle is a triangle with potentially shared edges
type ITriangle interface {
	Set(x1, y1, x2, y2, x3, y3 int)
	Draw(raster IRasterBuffer)
	Fill(raster IRasterBuffer)
}
