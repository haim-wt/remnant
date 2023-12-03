package physics

import (
	"gonum.org/v1/gonum/mat"
)

type RigidBody struct {
	Position     *mat.VecDense
	Velocity     *mat.VecDense
	Direction    *mat.VecDense
	Rotation     *mat.VecDense
	AngularVel   *mat.VecDense
	Acceleration *mat.VecDense
	Force        *mat.VecDense
	Torque       *mat.VecDense
	Mass         float64
}

func NewRigidBody(position *mat.VecDense) *RigidBody {
	return &RigidBody{
		Position:     mat.NewVecDense(3, []float64{0, 156, 0}),
		Velocity:     mat.NewVecDense(3, []float64{0, 0, 0}),
		Direction:    mat.NewVecDense(3, []float64{0, 0, 0}),
		Rotation:     mat.NewVecDense(3, []float64{0, 0, 0}),
		Acceleration: mat.NewVecDense(3, []float64{0, 0, 0}),
		AngularVel:   mat.NewVecDense(3, []float64{0, 0, 0}),
		Force:        mat.NewVecDense(3, []float64{0, 0, 0}),
		Torque:       mat.NewVecDense(3, []float64{0, 0, 0}),
		Mass:         5,
	}
}

func (rb *RigidBody) Update(dt float64) {

	rb.Acceleration.ScaleVec(1/rb.Mass, rb.Force)
	rb.Velocity.AddScaledVec(rb.Velocity, dt, rb.Acceleration)
	rb.Position.AddScaledVec(rb.Position, dt, rb.Velocity)

	rb.AngularVel.AddScaledVec(rb.AngularVel, dt, rb.Torque)
	rb.Direction.AddScaledVec(rb.Direction, dt, rb.AngularVel)

	rb.Force.ScaleVec(0, rb.Force)
	rb.Torque.ScaleVec(0, rb.Torque)
}

func (rb *RigidBody) ApplyForce(force *mat.VecDense) {
	rb.Force.AddVec(rb.Force, force)
}

func (rb *RigidBody) ApplyTorque(torque *mat.VecDense) {
	rb.Torque.AddVec(rb.Torque, torque)
}
