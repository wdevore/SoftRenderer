package smath

import (
	"fmt"
	"math"
)

// NewVector3 creates a Vector3 initialized to 0.0, 0.0, 0.0
func NewVector3() *Vector3 {
	v := new(Vector3)
	v.X = 0.0
	v.Y = 0.0
	v.Z = 0.0
	return v
}

// NewVector3With3Components creates a Vector3 initialized with x,y,z
func NewVector3With3Components(x, y, z float64) *Vector3 {
	v := new(Vector3)
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// NewVector3With2Components creates a Vector3 initialized with x,y and z = 0.0
func NewVector3With2Components(x, y float64) *Vector3 {
	v := new(Vector3)
	v.X = x
	v.Y = y
	v.Z = 0.0
	return v
}

// Clone returns a new copy this vector
func (v *Vector3) Clone() *Vector3 {
	c := new(Vector3)
	c.X = v.X
	c.Y = v.Y
	c.Z = v.Z
	return c
}

// Set3Components modifies x,y,z
func (v *Vector3) Set3Components(x, y, z float64) {
	v.X = x
	v.Y = y
	v.Z = z
}

// Set2Components modifies x,y only
func (v *Vector3) Set2Components(x, y float64) {
	v.X = x
	v.Y = y
}

// Set modifies x,y,z from source
func (v *Vector3) Set(source *Vector3) {
	v.X = source.X
	v.Y = source.Y
	v.Z = source.Z
}

// Add a Vector3 to this vector
func (v *Vector3) Add(src *Vector3) *Vector3 {
	v.X += src.X
	v.Y += src.Y
	v.Z += src.Z
	return v
}

// Add2Components adds x and y to this vector
func (v *Vector3) Add2Components(x, y float64) *Vector3 {
	v.X += x
	v.Y += y
	return v
}

// Sub subtracts a Vector3 to this vector
func (v *Vector3) Sub(src *Vector3) *Vector3 {
	v.X -= src.X
	v.Y -= src.Y
	v.Z -= src.Z
	return v
}

// Sub2Components subtracts x and y to this vector
func (v *Vector3) Sub2Components(x, y float64) *Vector3 {
	v.X -= x
	v.Y -= y
	return v
}

// ScaleBy scales this vector by s
func (v *Vector3) ScaleBy(s float64) *Vector3 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

// ScaleBy2Components scales this vector by sx and sy
func (v *Vector3) ScaleBy2Components(sx, sy float64) *Vector3 {
	v.X *= sx
	v.Y *= sy
	return v
}

// MulAdd scales and adds src to this vector
func (v *Vector3) MulAdd(src *Vector3, scalar float64) *Vector3 {
	v.X += src.X * scalar
	v.Y += src.Y * scalar
	v.Z += src.Z * scalar
	return v
}

// DivScalar 1/scales this vector
func (v *Vector3) DivScalar(scale float64) {
	v.X /= scale
	v.Y /= scale
	v.Z /= scale
}

// Length returns the euclidean length
func Length(x, y, z float64) float64 {
	return math.Sqrt(float64(x*x + y*y + z*z))
}

// Length returns the euclidean length
func (v *Vector3) Length() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z))
}

// LengthSquared returns the euclidean length squared
func LengthSquared(x, y, z float64) float64 {
	return x*x + y*y + z*z
}

// LengthSquared returns the euclidean length squared
func (v *Vector3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Normalize this instance
func (v *Vector3) Normalize() {
	l := v.Length()
	v.DivScalar(l)
}

// Equal makes an exact equality check. Use EqEpsilon, it is more realistic.
func (v *Vector3) Equal(other *Vector3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// EqEpsilon makes an approximate equality check. Preferred
func (v *Vector3) EqEpsilon(other *Vector3) bool {
	return (v.X-other.X) < Epsilon && (v.Y-other.Y) < Epsilon && (v.Z-other.Z) < Epsilon
}

// Distance finds the euclidean distance between the two specified vectors
func Distance(x1, y1, z1, x2, y2, z2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	c := z2 - z1

	return math.Sqrt(float64(a*a + b*b + c*c))
}

// Distance finds the euclidean distance between the two specified vectors
func (v *Vector3) Distance(src *Vector3) float64 {
	a := src.X - v.X
	b := src.Y - v.Y
	c := src.Z - v.Z

	return math.Sqrt(float64(a*a + b*b + c*c))
}

// DistanceSquared finds the euclidean distance between the two specified vectors squared
func DistanceSquared(x1, y1, z1, x2, y2, z2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	c := z2 - z1

	return a*a + b*b + c*c
}

// DistanceSquared finds the euclidean distance between the two specified vectors squared
func (v *Vector3) DistanceSquared(src *Vector3) float64 {
	a := src.X - v.X
	b := src.Y - v.Y
	c := src.Z - v.Z

	return a*a + b*b + c*c
}

// Dot returns the product between the two vectors
func Dot(x1, y1, z1, x2, y2, z2 float64) float64 {
	return x1*x2 + y1*y2 + z1*z2
}

// DotByComponent returns the product between the two vectors
func (v *Vector3) DotByComponent(x, y, z float64) float64 {
	return v.X*x + v.Y*y + v.Z*z
}

// Dot returns the product between the two vectors
func (v *Vector3) Dot(o *Vector3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// Cross sets this vector to the cross product between it and the other vector.
func (v *Vector3) Cross(o *Vector3) *Vector3 {
	v.Set3Components(
		v.Y*o.Z-v.Z*o.Y,
		v.Z*o.X-v.X*o.Z,
		v.X*o.Y-v.Y*o.X)

	return v
}

// --------------------------------------------------------------------------
// Transforms
// --------------------------------------------------------------------------

// Mul left-multiplies the vector by the given matrix, assuming the fourth (w) component of the vector is 1.
func (v *Vector3) Mul(m *Matrix4) *Vector3 {
	v.Set3Components(
		v.X*m.e[M00]+v.Y*m.e[M01]+v.Z*m.e[M02]+m.e[M03],
		v.X*m.e[M10]+v.Y*m.e[M11]+v.Z*m.e[M12]+m.e[M13],
		v.X*m.e[M20]+v.Y*m.e[M21]+v.Z*m.e[M22]+m.e[M23])
	return v
}

func (v Vector3) String() string {
	return fmt.Sprintf("<%f, %f, %f>", v.X, v.Y, v.Z)
}
