package smath

import "math"

// SetFromQuaternion axis/angle based on a quaternion.
func (aa *AxisAngle) SetFromQuaternion(q *Quaternion) {
	aa.Angle = 2 * math.Acos(q.W)
	s := math.Sin(aa.Angle / 2)
	aa.X = q.X / s
	aa.Y = q.Y / s
	aa.Z = q.Z / s
}
