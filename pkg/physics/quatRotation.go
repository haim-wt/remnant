package physics

import (
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
)

func CreateRotationQuaternion(angle float64, axis *mat.VecDense) quat.Number {
	s := math.Sin(angle / 2)
	return quat.Number{
		Real: math.Cos(angle / 2),
		Imag: axis.AtVec(0) * s,
		Jmag: axis.AtVec(1) * s,
		Kmag: axis.AtVec(2) * s,
	}
}

func RotateVectorByQuaternion(v *mat.VecDense, q quat.Number) *mat.VecDense {
	// Convert the vector to a quaternion (w=0)
	vecQuat := quat.Number{
		Imag: v.AtVec(0),
		Jmag: v.AtVec(1),
		Kmag: v.AtVec(2)}

	// Apply the rotation: q * vecQuat * conj(q)
	rotatedQuat := quat.Mul(quat.Mul(q, vecQuat), quat.Conj(q))

	// Extract the rotated vector
	return mat.NewVecDense(3, []float64{rotatedQuat.Imag, rotatedQuat.Jmag, rotatedQuat.Kmag})
}
