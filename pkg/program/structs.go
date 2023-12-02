package program

import (
	"gonum.org/v1/gonum/mat"
)

type Camera struct {
	Pos *mat.VecDense
	Dir *mat.VecDense
	Up  *mat.VecDense
	FOV float32
}

type Light struct {
	Position *mat.VecDense
}
