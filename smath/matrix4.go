package smath

import (
	"fmt"
	"math"
)

const (
	// M00 XX: Typically the unrotated X component for scaling, also the cosine of the
	// angle when rotated on the Y and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source X component
	// and added to the target X component.
	M00 = 0
	// M01 XY: Typically the negative sine of the angle when rotated on the Z axis.
	// On Vector3 multiplication this value is multiplied
	// with the source Y component and added to the target X component.
	M01 = 4
	// M02 XZ: Typically the sine of the angle when rotated on the Y axis.
	// On Vector3 multiplication this value is multiplied with the
	// source Z component and added to the target X component.
	M02 = 8
	// M03 XW: Typically the translation of the X component.
	// On Vector3 multiplication this value is added to the target X component.
	M03 = 12

	// M10 YX: Typically the sine of the angle when rotated on the Z axis.
	// On Vector3 multiplication this value is multiplied with the
	// source X component and added to the target Y component.
	M10 = 1
	// M11 YY: Typically the unrotated Y component for scaling, also the cosine
	// of the angle when rotated on the X and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source Y
	// component and added to the target Y component.
	M11 = 5
	// M12 YZ: Typically the negative sine of the angle when rotated on the X axis.
	// On Vector3 multiplication this value is multiplied
	// with the source Z component and added to the target Y component.
	M12 = 9
	// M13 YW: Typically the translation of the Y component.
	// On Vector3 multiplication this value is added to the target Y component.
	M13 = 13

	// M20 ZX: Typically the negative sine of the angle when rotated on the Y axis.
	// On Vector3 multiplication this value is multiplied
	// with the source X component and added to the target Z component.
	M20 = 2
	// M21 ZY: Typical the sine of the angle when rotated on the X axis.
	// On Vector3 multiplication this value is multiplied with the
	// source Y component and added to the target Z component.
	M21 = 6
	// M22 ZZ: Typically the unrotated Z component for scaling, also the cosine of the
	// angle when rotated on the X and/or Y axis.
	// On Vector3 multiplication this value is multiplied with the source Z component
	// and added to the target Z component.
	M22 = 10
	// M23 ZW: Typically the translation of the Z component.
	// On Vector3 multiplication this value is added to the target Z component.
	M23 = 14

	// M30 WX: Typically the value zero. On Vector3 multiplication this value is ignored.
	M30 = 3
	// M31 WY: Typically the value zero. On Vector3 multiplication this value is ignored.
	M31 = 7
	// M32 WZ: Typically the value zero. On Vector3 multiplication this value is ignored.
	M32 = 11
	// M33 WW: Typically the value one. On Vector3 multiplication this value is ignored.
	M33 = 15
)

// A temporary matrix for multiplication
var temp = NewMatrix4()

// NewMatrix4 creates a Matrix4 initialized to an identity matrix
func NewMatrix4() *Matrix4 {
	m := new(Matrix4)
	m.ToIdentity()
	return m
}

// --------------------------------------------------------------------------
// Translation
// --------------------------------------------------------------------------

// TranslateBy adds a translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *Matrix4) TranslateBy(v *Vector3) *Matrix4 {

	m.e[M03] += v.X
	m.e[M13] += v.Y
	m.e[M23] += v.Z

	return m
}

// TranslateBy3Comps adds a translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *Matrix4) TranslateBy3Comps(x, y, z float64) *Matrix4 {

	m.e[M03] += x
	m.e[M13] += y
	m.e[M23] += z

	return m
}

// TranslateBy2Comps adds a translational component to the matrix in the 4th column.
// Z is unmodified. The other columns are unmodified.
func (m *Matrix4) TranslateBy2Comps(x, y float64) *Matrix4 {

	m.e[M03] += x
	m.e[M13] += y

	return m
}

// SetTranslateByVector sets the translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *Matrix4) SetTranslateByVector(v *Vector3) *Matrix4 {
	m.ToIdentity()
	m.e[M03] = v.X
	m.e[M13] = v.Y
	m.e[M23] = v.Z

	return m
}

// SetTranslate3Comp sets the translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *Matrix4) SetTranslate3Comp(x, y, z float64) *Matrix4 {
	m.ToIdentity()

	m.e[M03] = x
	m.e[M13] = y
	m.e[M23] = z

	return m
}

// GetTranslation returns the translational components in 'out' Vector3 field.
func (m *Matrix4) GetTranslation(out *Vector3) {
	out.Set3Components(m.e[M03], m.e[M13], m.e[M23])
}

// --------------------------------------------------------------------------
// Rotation
// --------------------------------------------------------------------------

// SetRotation set a rotation matrix about Z axis. 'angle' is specified
// in radians.
//
//      [  M00  M01   _    _   ]
//      [  M10  M11   _    _   ]
//      [   _    _    _    _   ]
//      [   _    _    _    _   ]
func (m *Matrix4) SetRotation(angle float64) *Matrix4 {
	m.ToIdentity()

	if angle == 0 {
		return m
	}

	m.Rotation = angle

	// Column major
	c := math.Cos(angle)
	s := math.Sin(angle)

	m.e[M00] = c
	m.e[M01] = -s
	m.e[M10] = s
	m.e[M11] = c

	return m
}

// RotateBy postmultiplies this matrix with a (counter-clockwise) rotation matrix whose
// angle is specified in radians.
func (m *Matrix4) RotateBy(angle float64) *Matrix4 {
	if angle == 0.0 {
		return m
	}

	m.Rotation += angle

	// Column major
	c := math.Cos(angle)
	s := math.Sin(angle)

	temp.e[M00] = c
	temp.e[M01] = -s
	temp.e[M02] = 0.0
	temp.e[M03] = 0.0
	temp.e[M10] = s
	temp.e[M11] = c
	temp.e[M12] = 0.0
	temp.e[M13] = 0.0
	temp.e[M20] = 0.0
	temp.e[M21] = 0.0
	temp.e[M22] = 1.0
	temp.e[M23] = 0.0
	temp.e[M30] = 0.0
	temp.e[M31] = 0.0
	temp.e[M32] = 0.0
	temp.e[M33] = 1.0

	m.PostMultiply(temp)

	return m
}

// --------------------------------------------------------------------------
// Scale
// --------------------------------------------------------------------------

// ScaleBy scales the scale components.
func (m *Matrix4) ScaleBy(v *Vector3) *Matrix4 {
	m.Scale.Set(v)

	m.e[M00] *= v.X
	m.e[M11] *= v.Y
	m.e[M22] *= v.Z

	return m
}

// SetScale sets the scale components of an identity matrix and captures
// scale values into Scale property.
func (m *Matrix4) SetScale(v *Vector3) *Matrix4 {
	m.Scale.Set(v)

	m.ToIdentity()

	m.e[M00] = v.X
	m.e[M11] = v.Y
	m.e[M22] = v.Z

	return m
}

// SetScale3Comp sets the scale components of an identity matrix and captures
// scale values into Scale property.
func (m *Matrix4) SetScale3Comp(sx, sy, sz float64) *Matrix4 {
	m.Scale.Set3Components(sx, sy, sz)

	m.ToIdentity()

	m.e[M00] = sx
	m.e[M11] = sy
	m.e[M22] = sz

	return m
}

// SetScale2Comp sets the scale components of an identity matrix and captures
// scale values into Scale property where Z component = 1.0.
func (m *Matrix4) SetScale2Comp(sx, sy float64) *Matrix4 {
	m.Scale.Set3Components(sx, sy, 1.0)

	m.ToIdentity()

	m.e[M00] = sx
	m.e[M11] = sy
	m.e[M22] = 1.0

	return m
}

// GetScale returns the scale in 'out' field.
func (m *Matrix4) GetScale(out *Vector3) {
	out.Set(&m.Scale)
}

// PostScale postmultiplies this matrix with a scale matrix.
func (m *Matrix4) PostScale(sx, sy, sz float64) *Matrix4 {
	temp.e[M00] = sx
	temp.e[M01] = 0
	temp.e[M02] = 0
	temp.e[M03] = 0
	temp.e[M10] = 0
	temp.e[M11] = sy
	temp.e[M12] = 0
	temp.e[M13] = 0
	temp.e[M20] = 0
	temp.e[M21] = 0
	temp.e[M22] = sz
	temp.e[M23] = 0
	temp.e[M30] = 0
	temp.e[M31] = 0
	temp.e[M32] = 0
	temp.e[M33] = 1

	m.PostMultiply(temp)

	return m
}

// --------------------------------------------------------------------------
// Transforms
// --------------------------------------------------------------------------

// --------------------------------------------------------------------------
// Matrix methods
// --------------------------------------------------------------------------

// Multiply multiplies a * b and places result into 'out', (i.e. out = a * b)
func Multiply(a, b, out *Matrix4) {
	out.e[M00] = a.e[M00]*b.e[M00] + a.e[M01]*b.e[M10] + a.e[M02]*b.e[M20] + a.e[M03]*b.e[M30]
	out.e[M01] = a.e[M00]*b.e[M01] + a.e[M01]*b.e[M11] + a.e[M02]*b.e[M21] + a.e[M03]*b.e[M31]
	out.e[M02] = a.e[M00]*b.e[M02] + a.e[M01]*b.e[M12] + a.e[M02]*b.e[M22] + a.e[M03]*b.e[M32]
	out.e[M03] = a.e[M00]*b.e[M03] + a.e[M01]*b.e[M13] + a.e[M02]*b.e[M23] + a.e[M03]*b.e[M33]
	out.e[M10] = a.e[M10]*b.e[M00] + a.e[M11]*b.e[M10] + a.e[M12]*b.e[M20] + a.e[M13]*b.e[M30]
	out.e[M11] = a.e[M10]*b.e[M01] + a.e[M11]*b.e[M11] + a.e[M12]*b.e[M21] + a.e[M13]*b.e[M31]
	out.e[M12] = a.e[M10]*b.e[M02] + a.e[M11]*b.e[M12] + a.e[M12]*b.e[M22] + a.e[M13]*b.e[M32]
	out.e[M13] = a.e[M10]*b.e[M03] + a.e[M11]*b.e[M13] + a.e[M12]*b.e[M23] + a.e[M13]*b.e[M33]
	out.e[M20] = a.e[M20]*b.e[M00] + a.e[M21]*b.e[M10] + a.e[M22]*b.e[M20] + a.e[M23]*b.e[M30]
	out.e[M21] = a.e[M20]*b.e[M01] + a.e[M21]*b.e[M11] + a.e[M22]*b.e[M21] + a.e[M23]*b.e[M31]
	out.e[M22] = a.e[M20]*b.e[M02] + a.e[M21]*b.e[M12] + a.e[M22]*b.e[M22] + a.e[M23]*b.e[M32]
	out.e[M23] = a.e[M20]*b.e[M03] + a.e[M21]*b.e[M13] + a.e[M22]*b.e[M23] + a.e[M23]*b.e[M33]
	out.e[M30] = a.e[M30]*b.e[M00] + a.e[M31]*b.e[M10] + a.e[M32]*b.e[M20] + a.e[M33]*b.e[M30]
	out.e[M31] = a.e[M30]*b.e[M01] + a.e[M31]*b.e[M11] + a.e[M32]*b.e[M21] + a.e[M33]*b.e[M31]
	out.e[M32] = a.e[M30]*b.e[M02] + a.e[M31]*b.e[M12] + a.e[M32]*b.e[M22] + a.e[M33]*b.e[M32]
	out.e[M33] = a.e[M30]*b.e[M03] + a.e[M31]*b.e[M13] + a.e[M32]*b.e[M23] + a.e[M33]*b.e[M33]
}

// Multiply multiplies a * b and places result into this matrix, (i.e. this = a * b)
func (m *Matrix4) Multiply(a, b *Matrix4) {
	m.e[M00] = a.e[M00]*b.e[M00] + a.e[M01]*b.e[M10] + a.e[M02]*b.e[M20] + a.e[M03]*b.e[M30]
	m.e[M01] = a.e[M00]*b.e[M01] + a.e[M01]*b.e[M11] + a.e[M02]*b.e[M21] + a.e[M03]*b.e[M31]
	m.e[M02] = a.e[M00]*b.e[M02] + a.e[M01]*b.e[M12] + a.e[M02]*b.e[M22] + a.e[M03]*b.e[M32]
	m.e[M03] = a.e[M00]*b.e[M03] + a.e[M01]*b.e[M13] + a.e[M02]*b.e[M23] + a.e[M03]*b.e[M33]
	m.e[M10] = a.e[M10]*b.e[M00] + a.e[M11]*b.e[M10] + a.e[M12]*b.e[M20] + a.e[M13]*b.e[M30]
	m.e[M11] = a.e[M10]*b.e[M01] + a.e[M11]*b.e[M11] + a.e[M12]*b.e[M21] + a.e[M13]*b.e[M31]
	m.e[M12] = a.e[M10]*b.e[M02] + a.e[M11]*b.e[M12] + a.e[M12]*b.e[M22] + a.e[M13]*b.e[M32]
	m.e[M13] = a.e[M10]*b.e[M03] + a.e[M11]*b.e[M13] + a.e[M12]*b.e[M23] + a.e[M13]*b.e[M33]
	m.e[M20] = a.e[M20]*b.e[M00] + a.e[M21]*b.e[M10] + a.e[M22]*b.e[M20] + a.e[M23]*b.e[M30]
	m.e[M21] = a.e[M20]*b.e[M01] + a.e[M21]*b.e[M11] + a.e[M22]*b.e[M21] + a.e[M23]*b.e[M31]
	m.e[M22] = a.e[M20]*b.e[M02] + a.e[M21]*b.e[M12] + a.e[M22]*b.e[M22] + a.e[M23]*b.e[M32]
	m.e[M23] = a.e[M20]*b.e[M03] + a.e[M21]*b.e[M13] + a.e[M22]*b.e[M23] + a.e[M23]*b.e[M33]
	m.e[M30] = a.e[M30]*b.e[M00] + a.e[M31]*b.e[M10] + a.e[M32]*b.e[M20] + a.e[M33]*b.e[M30]
	m.e[M31] = a.e[M30]*b.e[M01] + a.e[M31]*b.e[M11] + a.e[M32]*b.e[M21] + a.e[M33]*b.e[M31]
	m.e[M32] = a.e[M30]*b.e[M02] + a.e[M31]*b.e[M12] + a.e[M32]*b.e[M22] + a.e[M33]*b.e[M32]
	m.e[M33] = a.e[M30]*b.e[M03] + a.e[M31]*b.e[M13] + a.e[M32]*b.e[M23] + a.e[M33]*b.e[M33]
}

// PreMultiply premultiplies 'b' matrix with 'this' and places the result into 'this' matrix.
// (i.e. this = b * this)
func (m *Matrix4) PreMultiply(b *Matrix4) {
	temp.e[M00] = b.e[M00]*m.e[M00] + b.e[M01]*m.e[M10] + b.e[M02]*m.e[M20] + b.e[M03]*m.e[M30]
	temp.e[M01] = b.e[M00]*m.e[M01] + b.e[M01]*m.e[M11] + b.e[M02]*m.e[M21] + b.e[M03]*m.e[M31]
	temp.e[M02] = b.e[M00]*m.e[M02] + b.e[M01]*m.e[M12] + b.e[M02]*m.e[M22] + b.e[M03]*m.e[M32]
	temp.e[M03] = b.e[M00]*m.e[M03] + b.e[M01]*m.e[M13] + b.e[M02]*m.e[M23] + b.e[M03]*m.e[M33]
	temp.e[M10] = b.e[M10]*m.e[M00] + b.e[M11]*m.e[M10] + b.e[M12]*m.e[M20] + b.e[M13]*m.e[M30]
	temp.e[M11] = b.e[M10]*m.e[M01] + b.e[M11]*m.e[M11] + b.e[M12]*m.e[M21] + b.e[M13]*m.e[M31]
	temp.e[M12] = b.e[M10]*m.e[M02] + b.e[M11]*m.e[M12] + b.e[M12]*m.e[M22] + b.e[M13]*m.e[M32]
	temp.e[M13] = b.e[M10]*m.e[M03] + b.e[M11]*m.e[M13] + b.e[M12]*m.e[M23] + b.e[M13]*m.e[M33]
	temp.e[M20] = b.e[M20]*m.e[M00] + b.e[M21]*m.e[M10] + b.e[M22]*m.e[M20] + b.e[M23]*m.e[M30]
	temp.e[M21] = b.e[M20]*m.e[M01] + b.e[M21]*m.e[M11] + b.e[M22]*m.e[M21] + b.e[M23]*m.e[M31]
	temp.e[M22] = b.e[M20]*m.e[M02] + b.e[M21]*m.e[M12] + b.e[M22]*m.e[M22] + b.e[M23]*m.e[M32]
	temp.e[M23] = b.e[M20]*m.e[M03] + b.e[M21]*m.e[M13] + b.e[M22]*m.e[M23] + b.e[M23]*m.e[M33]
	temp.e[M30] = b.e[M30]*m.e[M00] + b.e[M31]*m.e[M10] + b.e[M32]*m.e[M20] + b.e[M33]*m.e[M30]
	temp.e[M31] = b.e[M30]*m.e[M01] + b.e[M31]*m.e[M11] + b.e[M32]*m.e[M21] + b.e[M33]*m.e[M31]
	temp.e[M32] = b.e[M30]*m.e[M02] + b.e[M31]*m.e[M12] + b.e[M32]*m.e[M22] + b.e[M33]*m.e[M32]
	temp.e[M33] = b.e[M30]*m.e[M03] + b.e[M31]*m.e[M13] + b.e[M32]*m.e[M23] + b.e[M33]*m.e[M33]

	// Place results in "this"
	m.e[M00] = temp.e[M00]
	m.e[M01] = temp.e[M01]
	m.e[M02] = temp.e[M02]
	m.e[M03] = temp.e[M03]
	m.e[M10] = temp.e[M10]
	m.e[M11] = temp.e[M11]
	m.e[M12] = temp.e[M12]
	m.e[M13] = temp.e[M13]
	m.e[M20] = temp.e[M20]
	m.e[M21] = temp.e[M21]
	m.e[M22] = temp.e[M22]
	m.e[M23] = temp.e[M23]
	m.e[M30] = temp.e[M30]
	m.e[M31] = temp.e[M31]
	m.e[M32] = temp.e[M32]
	m.e[M33] = temp.e[M33]
}

// PostMultiply postmultiplies 'b' matrix with 'this' and places the result into 'this' matrix.
// (i.e. this = this * b)
func (m *Matrix4) PostMultiply(b *Matrix4) {
	temp.e[M00] = b.e[M00]*m.e[M00] + b.e[M01]*m.e[M10] + b.e[M02]*m.e[M20] + b.e[M03]*m.e[M30]
	temp.e[M01] = b.e[M00]*m.e[M01] + b.e[M01]*m.e[M11] + b.e[M02]*m.e[M21] + b.e[M03]*m.e[M31]
	temp.e[M02] = b.e[M00]*m.e[M02] + b.e[M01]*m.e[M12] + b.e[M02]*m.e[M22] + b.e[M03]*m.e[M32]
	temp.e[M03] = b.e[M00]*m.e[M03] + b.e[M01]*m.e[M13] + b.e[M02]*m.e[M23] + b.e[M03]*m.e[M33]
	temp.e[M10] = b.e[M10]*m.e[M00] + b.e[M11]*m.e[M10] + b.e[M12]*m.e[M20] + b.e[M13]*m.e[M30]
	temp.e[M11] = b.e[M10]*m.e[M01] + b.e[M11]*m.e[M11] + b.e[M12]*m.e[M21] + b.e[M13]*m.e[M31]
	temp.e[M12] = b.e[M10]*m.e[M02] + b.e[M11]*m.e[M12] + b.e[M12]*m.e[M22] + b.e[M13]*m.e[M32]
	temp.e[M13] = b.e[M10]*m.e[M03] + b.e[M11]*m.e[M13] + b.e[M12]*m.e[M23] + b.e[M13]*m.e[M33]
	temp.e[M20] = b.e[M20]*m.e[M00] + b.e[M21]*m.e[M10] + b.e[M22]*m.e[M20] + b.e[M23]*m.e[M30]
	temp.e[M21] = b.e[M20]*m.e[M01] + b.e[M21]*m.e[M11] + b.e[M22]*m.e[M21] + b.e[M23]*m.e[M31]
	temp.e[M22] = b.e[M20]*m.e[M02] + b.e[M21]*m.e[M12] + b.e[M22]*m.e[M22] + b.e[M23]*m.e[M32]
	temp.e[M23] = b.e[M20]*m.e[M03] + b.e[M21]*m.e[M13] + b.e[M22]*m.e[M23] + b.e[M23]*m.e[M33]
	temp.e[M30] = b.e[M30]*m.e[M00] + b.e[M31]*m.e[M10] + b.e[M32]*m.e[M20] + b.e[M33]*m.e[M30]
	temp.e[M31] = b.e[M30]*m.e[M01] + b.e[M31]*m.e[M11] + b.e[M32]*m.e[M21] + b.e[M33]*m.e[M31]
	temp.e[M32] = b.e[M30]*m.e[M02] + b.e[M31]*m.e[M12] + b.e[M32]*m.e[M22] + b.e[M33]*m.e[M32]
	temp.e[M33] = b.e[M30]*m.e[M03] + b.e[M31]*m.e[M13] + b.e[M32]*m.e[M23] + b.e[M33]*m.e[M33]

	// Place results in "this"
	m.e[M00] = temp.e[M00]
	m.e[M01] = temp.e[M01]
	m.e[M02] = temp.e[M02]
	m.e[M03] = temp.e[M03]
	m.e[M10] = temp.e[M10]
	m.e[M11] = temp.e[M11]
	m.e[M12] = temp.e[M12]
	m.e[M13] = temp.e[M13]
	m.e[M20] = temp.e[M20]
	m.e[M21] = temp.e[M21]
	m.e[M22] = temp.e[M22]
	m.e[M23] = temp.e[M23]
	m.e[M30] = temp.e[M30]
	m.e[M31] = temp.e[M31]
	m.e[M32] = temp.e[M32]
	m.e[M33] = temp.e[M33]
}

// MultiplyIntoA multiplies a * b and places result into 'a', (i.e. a = a * b)
func MultiplyIntoA(a, b *Matrix4) {
	temp.e[M00] = a.e[M00]*b.e[M00] + a.e[M01]*b.e[M10] + a.e[M02]*b.e[M20] + a.e[M03]*b.e[M30]
	temp.e[M01] = a.e[M00]*b.e[M01] + a.e[M01]*b.e[M11] + a.e[M02]*b.e[M21] + a.e[M03]*b.e[M31]
	temp.e[M02] = a.e[M00]*b.e[M02] + a.e[M01]*b.e[M12] + a.e[M02]*b.e[M22] + a.e[M03]*b.e[M32]
	temp.e[M03] = a.e[M00]*b.e[M03] + a.e[M01]*b.e[M13] + a.e[M02]*b.e[M23] + a.e[M03]*b.e[M33]
	temp.e[M10] = a.e[M10]*b.e[M00] + a.e[M11]*b.e[M10] + a.e[M12]*b.e[M20] + a.e[M13]*b.e[M30]
	temp.e[M11] = a.e[M10]*b.e[M01] + a.e[M11]*b.e[M11] + a.e[M12]*b.e[M21] + a.e[M13]*b.e[M31]
	temp.e[M12] = a.e[M10]*b.e[M02] + a.e[M11]*b.e[M12] + a.e[M12]*b.e[M22] + a.e[M13]*b.e[M32]
	temp.e[M13] = a.e[M10]*b.e[M03] + a.e[M11]*b.e[M13] + a.e[M12]*b.e[M23] + a.e[M13]*b.e[M33]
	temp.e[M20] = a.e[M20]*b.e[M00] + a.e[M21]*b.e[M10] + a.e[M22]*b.e[M20] + a.e[M23]*b.e[M30]
	temp.e[M21] = a.e[M20]*b.e[M01] + a.e[M21]*b.e[M11] + a.e[M22]*b.e[M21] + a.e[M23]*b.e[M31]
	temp.e[M22] = a.e[M20]*b.e[M02] + a.e[M21]*b.e[M12] + a.e[M22]*b.e[M22] + a.e[M23]*b.e[M32]
	temp.e[M23] = a.e[M20]*b.e[M03] + a.e[M21]*b.e[M13] + a.e[M22]*b.e[M23] + a.e[M23]*b.e[M33]
	temp.e[M30] = a.e[M30]*b.e[M00] + a.e[M31]*b.e[M10] + a.e[M32]*b.e[M20] + a.e[M33]*b.e[M30]
	temp.e[M31] = a.e[M30]*b.e[M01] + a.e[M31]*b.e[M11] + a.e[M32]*b.e[M21] + a.e[M33]*b.e[M31]
	temp.e[M32] = a.e[M30]*b.e[M02] + a.e[M31]*b.e[M12] + a.e[M32]*b.e[M22] + a.e[M33]*b.e[M32]
	temp.e[M33] = a.e[M30]*b.e[M03] + a.e[M31]*b.e[M13] + a.e[M32]*b.e[M23] + a.e[M33]*b.e[M33]

	a.e[M00] = temp.e[M00]
	a.e[M01] = temp.e[M01]
	a.e[M02] = temp.e[M02]
	a.e[M03] = temp.e[M03]
	a.e[M10] = temp.e[M10]
	a.e[M11] = temp.e[M11]
	a.e[M12] = temp.e[M12]
	a.e[M13] = temp.e[M13]
	a.e[M20] = temp.e[M20]
	a.e[M21] = temp.e[M21]
	a.e[M22] = temp.e[M22]
	a.e[M23] = temp.e[M23]
	a.e[M30] = temp.e[M30]
	a.e[M31] = temp.e[M31]
	a.e[M32] = temp.e[M32]
	a.e[M33] = temp.e[M33]
}

// PostTranslate postmultiplies this matrix by a translation matrix.
// Postmultiplication is also used by OpenGL ES.
func (m *Matrix4) PostTranslate(tx, ty, tz float64) *Matrix4 {
	temp.e[M00] = 1.0
	temp.e[M01] = 0.0
	temp.e[M02] = 0.0
	temp.e[M03] = tx
	temp.e[M10] = 0.0
	temp.e[M11] = 1.0
	temp.e[M12] = 0.0
	temp.e[M13] = ty
	temp.e[M20] = 0.0
	temp.e[M21] = 0.0
	temp.e[M22] = 1.0
	temp.e[M23] = tz
	temp.e[M30] = 0.0
	temp.e[M31] = 0.0
	temp.e[M32] = 0.0
	temp.e[M33] = 1.0

	// this = this * temp;
	m.PostMultiply(temp)

	return m
}

// --------------------------------------------------------------------------
// Projections
// --------------------------------------------------------------------------

// SetToOrtho sets the matrix for a 2d ortho graphic projection
func (m *Matrix4) SetToOrtho(left, right, bottom, top, near, far float64) *Matrix4 {
	m.ToIdentity()

	xorth := 2.0 / (right - left)
	yorth := 2.0 / (top - bottom)
	zorth := -2.0 / (far - near)

	tx := -(right + left) / (right - left)
	ty := -(top + bottom) / (top - bottom)
	tz := -(far + near) / (far - near)

	m.e[M00] = xorth
	m.e[M10] = 0.0
	m.e[M20] = 0.0
	m.e[M30] = 0.0
	m.e[M01] = 0.0
	m.e[M11] = yorth
	m.e[M21] = 0.0
	m.e[M31] = 0.0
	m.e[M02] = 0.0
	m.e[M12] = 0.0
	m.e[M22] = zorth
	m.e[M32] = 0.0
	m.e[M03] = tx
	m.e[M13] = ty
	m.e[M23] = tz
	m.e[M33] = 1.0

	return m
}

// --------------------------------------------------------------------------
// Misc
// --------------------------------------------------------------------------

// C returns a cell value based on Mxx index
func (m *Matrix4) C(i int) float64 {
	return m.e[i]
}

// Clone returns a clone of this matrix
func (m *Matrix4) Clone() *Matrix4 {
	c := new(Matrix4)
	c.Set(m)
	return c
}

// Set copies src into this matrix
func (m *Matrix4) Set(src *Matrix4) *Matrix4 {
	m.e[M00] = src.e[M00]
	m.e[M01] = src.e[M01]
	m.e[M02] = src.e[M02]
	m.e[M03] = src.e[M03]

	m.e[M10] = src.e[M10]
	m.e[M11] = src.e[M11]
	m.e[M12] = src.e[M12]
	m.e[M13] = src.e[M13]

	m.e[M20] = src.e[M20]
	m.e[M21] = src.e[M21]
	m.e[M22] = src.e[M22]
	m.e[M23] = src.e[M23]

	m.e[M30] = src.e[M30]
	m.e[M31] = src.e[M31]
	m.e[M32] = src.e[M32]
	m.e[M33] = src.e[M33]

	return m
}

// ToIdentity set this matrix to the identity matrix
func (m *Matrix4) ToIdentity() {
	m.e[M00] = 1.0
	m.e[M01] = 0.0
	m.e[M02] = 0.0
	m.e[M03] = 0.0

	m.e[M10] = 0.0
	m.e[M11] = 1.0
	m.e[M12] = 0.0
	m.e[M13] = 0.0

	m.e[M20] = 0.0
	m.e[M21] = 0.0
	m.e[M22] = 1.0
	m.e[M23] = 0.0

	m.e[M30] = 0.0
	m.e[M31] = 0.0
	m.e[M32] = 0.0
	m.e[M33] = 1.0
}

func (m Matrix4) String() string {
	s := fmt.Sprintf("[%f, %f, %f, %f]\n", m.e[M00], m.e[M01], m.e[M02], m.e[M03])
	s += fmt.Sprintf("[%f, %f, %f, %f]\n", m.e[M10], m.e[M11], m.e[M12], m.e[M13])
	s += fmt.Sprintf("[%f, %f, %f, %f]\n", m.e[M20], m.e[M21], m.e[M22], m.e[M23])
	s += fmt.Sprintf("[%f, %f, %f, %f]", m.e[M30], m.e[M31], m.e[M32], m.e[M33])
	return s
}
