package ship

import (
	"remnant/pkg/input"
	"remnant/pkg/physics"

	"github.com/go-gl/glfw/v3.3/glfw"
	"gonum.org/v1/gonum/mat"
)

type Movement struct {
	Forward   *input.Key
	Backward  *input.Key
	Left      *input.Key
	Right     *input.Key
	Up        *input.Key
	Down      *input.Key
	RollRight *input.Key
	RollLeft  *input.Key
}

func (m *Movement) UpdateMovement(window *glfw.Window, forward *mat.VecDense, up *mat.VecDense) (*mat.VecDense, float64) {
	d := mat.NewVecDense(3, []float64{0, 0, 0})
	d.CopyVec(forward)

	direction := mat.NewVecDense(3, []float64{0, 0, 0})
	if m.Forward.UpdateKeyState(window) {
		direction.AddVec(direction, d)
	}
	d.ScaleVec(-1, d)
	if m.Backward.UpdateKeyState(window) {
		direction.AddVec(direction, d)
	}

	d.CopyVec(up)
	if m.Up.UpdateKeyState(window) {
		direction.AddVec(direction, d)
	}
	d.ScaleVec(-1, d)
	if m.Down.UpdateKeyState(window) {
		direction.AddVec(direction, d)
	}

	d = physics.Cross(forward, up)
	if m.Left.UpdateKeyState(window) {
		direction.AddVec(direction, d)
	}

	d.ScaleVec(-1, d)
	if m.Right.UpdateKeyState(window) {
		direction.AddVec(direction, d)
	}

	roll := 0.0
	if m.RollLeft.UpdateKeyState(window) {
		roll += 0.02
	}

	if m.RollRight.UpdateKeyState(window) {
		roll -= 0.02
	}

	norm := direction.Norm(2)
	if norm > 0 {
		direction.ScaleVec(1/direction.Norm(2), direction)
	}
	return direction, roll
}

type Ship struct {
	*physics.RigidBody
	Movement *Movement
}

func NewShip(position *mat.VecDense) *Ship {
	return &Ship{
		RigidBody: physics.NewRigidBody(position),
		Movement:  &Movement{},
	}
}
