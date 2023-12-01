package main

import "gonum.org/v1/gonum/mat"

type Scene struct {
	Light   Light
	Camera  Camera
	Program Program
}

func NewScene() *Scene {
	return &Scene{
		Light: Light{
			Position: mat.NewVecDense(3, []float64{0, 64, -64}),
		},
		Camera: Camera{
			Pos: mat.NewVecDense(3, []float64{-32, 0, -32}),
			Dir: mat.NewVecDense(3, []float64{0, 0, 1}),
			Up:  mat.NewVecDense(3, []float64{0, 1, 0}),
			FOV: float32(60),
		},
		Program: Program{},
	}
}
