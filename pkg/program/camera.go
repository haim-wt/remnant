package program

import (
	"remnant/pkg/physics"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
)

type Camera struct {
	Pos *mat.VecDense
	Dir *mat.VecDense
	Up  *mat.VecDense
	FOV float32
}

func NewCamera(positon *mat.VecDense, fov float32) *Camera {
	return &Camera{
		Pos: positon,
		Dir: mat.NewVecDense(3, []float64{0, 0, 1}),
		Up:  mat.NewVecDense(3, []float64{0, 1, 0}),
		FOV: fov,
	}
}

func (c *Camera) Rotate(xRad, yRad float64) {

	qx := physics.CreateRotationQuaternion(xRad, c.Up)
	qy := physics.CreateRotationQuaternion(yRad, physics.Cross(c.Up, c.Dir))

	combinedRotation := quat.Mul(qx, qy)
	rotatedVec := physics.RotateVectorByQuaternion(c.Dir, combinedRotation)
	rotatedUp := physics.RotateVectorByQuaternion(c.Up, combinedRotation)

	c.Dir = rotatedVec
	c.Up = rotatedUp
}

func (c *Camera) RotateZ(xRad float64) {
	qz := physics.CreateRotationQuaternion(xRad, c.Dir)

	rotatedVec := physics.RotateVectorByQuaternion(c.Dir, qz)
	rotatedUp := physics.RotateVectorByQuaternion(c.Up, qz)

	c.Dir = rotatedVec
	c.Up = rotatedUp
}
