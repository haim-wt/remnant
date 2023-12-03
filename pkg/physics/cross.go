package physics

import "gonum.org/v1/gonum/mat"

func Cross(a *mat.VecDense, b *mat.VecDense) *mat.VecDense {
	ax := float64(a.AtVec(0))
	ay := float64(a.AtVec(1))
	az := float64(a.AtVec(2))

	bx := float64(b.AtVec(0))
	by := float64(b.AtVec(1))
	bz := float64(b.AtVec(2))

	return mat.NewVecDense(3, []float64{
		ay*bz - az*by,
		az*bx - ax*bz,
		ax*by - ay*bx,
	})
}
