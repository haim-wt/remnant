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

func (c *Camera) Rotate(xRad, yRad float64) {
	qx := physics.CreateRotationQuaternion(xRad, mat.NewVecDense(3, []float64{0, 1, 0}))
	qy := physics.CreateRotationQuaternion(yRad, mat.NewVecDense(3, []float64{1, 0, 0}))

	rotatedVec := physics.RotateVectorByQuaternion(c.Dir, qy)
	rotatedVec = physics.RotateVectorByQuaternion(rotatedVec, qx)

	combinedRotation := quat.Mul(qy, qx)
	rotatedUp := physics.RotateVectorByQuaternion(c.Up, combinedRotation)

	c.Dir = rotatedVec
	c.Up = rotatedUp
}
