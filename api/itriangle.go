package api

// ITriangle is a triangle with potentially shared edges
type ITriangle interface {
	Set(x1, y1, x2, y2, x3, y3 int)
	SetWithZ(x1, y1 int, z1 float32, x2, y2 int, z2 float32, x3, y3 int, z3 float32)
	Draw(raster IRasterBuffer)
	Fill(raster IRasterBuffer)
}
