package program

import (
	"gonum.org/v1/gonum/mat"
)

type Light struct {
	Position *mat.VecDense
}

func NewLight(direction *mat.VecDense) *Light {
	return &Light{
		Position: direction,
	}
}
