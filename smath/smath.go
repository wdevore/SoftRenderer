// Package smath contains:
// general math functions, constants, matrices and vectors
package smath

import (
	"math"
)

// Note: Special character symbols are entered by first using Ctrl+Shift+u 'together'
// and then releasing. Then enter a code, for example, Pi is '03c0'.
// Finally, hit "enter" key.

// Vector3 contains base components
type Vector3 struct {
	X, Y, Z float64
}

// Matrix4 represents a column major opengl array.
type Matrix4 struct {
	e [16]float64

	// Rotation is in radians
	Rotation float64
	Scale    Vector3
}

// Quaternion represents a quaternion W+X*i+Y*j+Z*k
type Quaternion struct {
	W float64 // Scalar component
	X float64 // i component
	Y float64 // j component
	Z float64 // k component
}

// AxisAngle represents an axis and angle
type AxisAngle struct {
	X     float64 // i component
	Y     float64 // j component
	Z     float64 // k component
	Angle float64
}

// Rectangle is represented in a coordinate system defined and shaped as follows:
//    Left,Top
//        ^---------.
//        |         |
// Y axis |         |
//        |         |
//        |         |
//        .---------> Right, Bottom
//          X axis
// This orientation is configured via an orthographic projection
type Rectangle struct {
	// Upper/Lower Components used mostly for rendering
	Top, Left, Bottom, Right float32
	// Components typically used for computations
	Width, Height float32

	Centered bool
}

// PI or π = ~3.1415927
const PI = 3.1415927

// PI2 or 2π
const PI2 = PI * 2

// multiply by this to convert from radians to degrees
const radiansToDegrees = 180.0 / PI

// multiply by this to convert from degrees to radians
const degreesToRadians = PI / 180.0

// Epsilon = 0.00001
const Epsilon = 0.00001 // ~32 bits

// ToRadians converts degrees to radians
func ToRadians(degrees float32) float32 {
	return degrees * degreesToRadians
}

// ToDegrees converts radians to degrees
func ToDegrees(radians float32) float32 {
	return radians * radiansToDegrees
}

// IsEqual compares to floats based on an Epsilon float
func IsEqual(a float32, b float32) bool {
	return math.Abs(float64(a-b)) <= Epsilon
}

// Min32 return the float32 min value
func Min32(a, b float32) float32 {
	if a <= b {
		return a
	}
	return b
}

// Max32 return the float32 max value
func Max32(a, b float32) float32 {
	if a >= b {
		return a
	}
	return b
}
