package graphics

import "SoftRenderer/smath"

const (
	LGNSegs = 4
	NSEGS   = 1 << LGNSegs
)

// ArcBall is a basic camera control
type ArcBall struct {
	Radius float64

	qNow  smath.Quaternion
	qDown smath.Quaternion
	qDrag smath.Quaternion

	vNow   smath.Vector3
	vDown  smath.Vector3
	vFrom  smath.Vector3
	vTo    smath.Vector3
	vrFrom smath.Vector3
	vrTo   smath.Vector3
	mNow   smath.Matrix4
	mDown  smath.Matrix4

	dragging bool

	vBallMouse    smath.Vector3
	q             smath.Quaternion
	qConj         smath.Quaternion
	pts           [NSEGS + 1]*smath.Vector3
	vVector       smath.Vector3
	vBase         smath.Vector3
	vDirection    smath.Vector3
	pCanvasSize   smath.Vector3
	vCanvasCenter smath.Vector3
	vLos          smath.Vector3

	// private AxisAngle4f aaYaxis = new AxisAngle4f(0.0f, 1.0f, 0.0f, 180.0f/57.3f);

	// Position represented as a matrix
	translation smath.Matrix4

	/*
	 * Orientation is expressed as a series of rotations
	 */
	rotation smath.Matrix4
	/*
	 * The combined rotation and translation.
	 */
	affineTransform smath.Matrix4

	p1 smath.Vector3
	p2 smath.Vector3

	/**
	 * Ken's arcball must be based on 2D screen coords where the +y is upwards respectively.
	 * However, SWT 2D screen coords are such that +y is downward.
	 * In order to sync with Ken's arcball I flip the screen coords to match Ken's
	 * y = (Height - y).
	 * Also, a conjugate is needed as well as a (-y).
	 */
	screenYOrientation bool // default to GL, Y is upwards.
}

// NewArcBall creates Ken's arc ball
func NewArcBall() *ArcBall {
	ab := new(ArcBall)
	return ab
}

// Place sets the center and size of the controller.
func (ab *ArcBall) Place(v *smath.Vector3, r float64) {
	ab.vCanvasCenter.Set(v)
	ab.Radius = r
}

func (ab *ArcBall) remapScreenCoords(p *smath.Vector3) {
	if !ab.screenYOrientation {
		ab.vNow.Set3Components(p.X, ab.pCanvasSize.Y-p.Y, 0.0)
	} else {
		ab.vNow.Set3Components(p.X, p.Y, 0.0)
	}
}

// Mouse incorporatea the given mouse position.
func (ab *ArcBall) Mouse(mp *smath.Vector3) {
	ab.remapScreenCoords(mp)
}

// Update updates the arcball
func (ab *ArcBall) Update() {
	// ab.vTo.Set(ab.mouseOnSphere(ab.vNow))

	// if ab.dragging {
	// 	qDrag = ab.mapFromBallPoints(ab.vFrom, ab.vTo)
	// 	ab.qNow.mul(ab.qDrag, ab.qDown)
	// }

	// ab.mapToBallPoints(ab.qDown, ab.vrFrom, ab.vrTo)
	// ab.q.Set(&ab.qNow)
}

// GetMatrix returns the ball's equivalent matrix
func (ab *ArcBall) GetMatrix() *smath.Matrix4 {
	return &ab.mNow
}

// GetRotationFromAxisAngle gets the arcball's rotation in radians.
func (ab *ArcBall) GetRotationFromAxisAngle(aa *smath.AxisAngle) {
	aa.SetFromQuaternion(&ab.qNow)
	ab.computeLOS()
}

func (ab *ArcBall) computeLOS() {
	ab.qConj.Set(&ab.qNow)
	ab.qConj.Conjugate()
	ab.q.SetFromComponents(0.0, 0.0, -1.0, 0.0) // -Z
	// qp := smath.Prod(ab.qConj, ab.q)
	qp := smath.Prod(ab.qNow, ab.qConj, ab.q)
	ab.q.Set(&qp)
	// ab.q.mul(ab.qNow)

	ab.vLos.Set3Components(ab.q.X, ab.q.Y, ab.q.Z) // the camera's los
	ab.vLos.Normalize()
}
