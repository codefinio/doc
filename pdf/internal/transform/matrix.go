

package transform

import (
	"fmt"
	"math"

	"github.com/codefinio/doc/common"
)

// Matrix is a linear transform matrix in homogenous coordinates.
// PDF coordinate transforms are always affine so we only need 6 of these. See newMatrix.
type Matrix [9]float64

// IdentityMatrix returns the identity transform.
func IdentityMatrix() Matrix {
	return NewMatrix(1, 0, 0, 1, 0, 0)
}

// TranslationMatrix returns a matrix that translates by `tx`, `ty`.
func TranslationMatrix(tx, ty float64) Matrix {
	return NewMatrix(1, 0, 0, 1, tx, ty)
}

// NewMatrix returns an affine transform matrix laid out in homogenous coordinates as
//      a  b  0
//      c  d  0
//      tx ty 1
func NewMatrix(a, b, c, d, tx, ty float64) Matrix {
	m := Matrix{
		a, b, 0,
		c, d, 0,
		tx, ty, 1,
	}
	m.clampRange()
	return m
}

// String returns a string describing `m`.
func (m Matrix) String() string {
	a, b, c, d, tx, ty := m[0], m[1], m[3], m[4], m[6], m[7]
	return fmt.Sprintf("[%.4f,%.4f,%.4f,%.4f:%.4f,%.4f]", a, b, c, d, tx, ty)
}

// Set sets `m` to affine transform a,b,c,d,tx,ty.
func (m *Matrix) Set(a, b, c, d, tx, ty float64) {
	m[0], m[1] = a, b
	m[3], m[4] = c, d
	m[6], m[7] = tx, ty
	m.clampRange()
}

// Concat sets `m` to `b` × `m`.
// `b` needs to be created by newMatrix. i.e. It must be an affine transform.
//    m00 m01 0     b00 b01 0     m00*b00 + m01*b01        m00*b10 + m01*b11        0
//    m10 m11 0  ×  b10 b11 0  =  m10*b00 + m11*b01        m10*b10 + m11*b11        0
//    m20 m21 1     b20 b21 1     m20*b00 + m21*b10 + b20  m20*b01 + m21*b11 + b21  1
func (m *Matrix) Concat(b Matrix) {
	*m = Matrix{
		b[0]*m[0] + b[1]*m[3], b[0]*m[1] + b[1]*m[4], 0,
		b[3]*m[0] + b[4]*m[3], b[3]*m[1] + b[4]*m[4], 0,
		b[6]*m[0] + b[7]*m[3] + m[6], b[6]*m[1] + b[7]*m[4] + m[7], 1,
	}
	m.clampRange()
}

// Mult returns `b` × `m`.
func (m Matrix) Mult(b Matrix) Matrix {
	m.Concat(b)
	return m
}

// Translate appends a translation of `dx`,`dy` to `m`.
// m.Translate(dx, dy) is equivalent to m.Concat(NewMatrix(1, 0, 0, 1, dx, dy))
func (m *Matrix) Translate(dx, dy float64) {
	m[6] += dx
	m[7] += dy
	m.clampRange()
}

// Translation returns the translation part of `m`.
func (m *Matrix) Translation() (float64, float64) {
	return m[6], m[7]
}

// Transform returns coordinates `x`,`y` transformed by `m`.
func (m *Matrix) Transform(x, y float64) (float64, float64) {
	xp := x*m[0] + y*m[1] + m[6]
	yp := x*m[3] + y*m[4] + m[7]
	return xp, yp
}

// ScalingFactorX returns the X scaling of the affine transform.
func (m *Matrix) ScalingFactorX() float64 {
	return math.Hypot(m[0], m[1])
}

// ScalingFactorY returns the Y scaling of the affine transform.
func (m *Matrix) ScalingFactorY() float64 {
	return math.Hypot(m[3], m[4])
}

// Angle returns the angle of the affine transform in `m` in degrees.
func (m *Matrix) Angle() float64 {
	theta := math.Atan2(-m[1], m[0])
	if theta < 0.0 {
		theta += 2 * math.Pi
	}
	return theta / math.Pi * 180.0

}

// clampRange forces `m` to have reasonable values. It is a guard against crazy values in corrupt PDF files.
// Currently it clamps elements to [-maxAbsNumber, -maxAbsNumber] to avoid floating point exceptions.
func (m *Matrix) clampRange() {
	for i, x := range m {
		if x > maxAbsNumber {
			common.Log.Debug("CLAMP: %d -> %d", x, maxAbsNumber)
			m[i] = maxAbsNumber
		} else if x < -maxAbsNumber {
			common.Log.Debug("CLAMP: %d -> %d", x, -maxAbsNumber)
			m[i] = -maxAbsNumber
		}
	}
}

// maxAbsNumber defines the maximum absolute value of allowed practical matrix element values as needed
// to avoid floating point exceptions.
// TODO(gunnsth): Add reference or point to a specific example PDF that validates this.
const maxAbsNumber = 1e9
