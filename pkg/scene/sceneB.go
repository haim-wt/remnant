package scene

import (
	"fmt"
	"math/rand"
	"remnant/internal/controller"
	"remnant/pkg/input"
	"remnant/pkg/program"
	"remnant/pkg/ship"

	"github.com/go-gl/glfw/v3.3/glfw"
	"gonum.org/v1/gonum/mat"
)

type SceneB struct {
	*controller.Controller
	camera *program.Camera
	light  *program.Light
	person *ship.Ship
}

func NewSceneB(ctr *controller.Controller) *SceneB {
	sceneB := &SceneB{
		Controller: ctr,
		light:      program.NewLight(mat.NewVecDense(3, []float64{100, 100, 0})),
		camera:     program.NewCamera(mat.NewVecDense(3, []float64{0, 0, -16}), 60),
		person:     ship.NewShip(mat.NewVecDense(3, []float64{0, 0, -16})),
	}

	sceneB.person.Movement = &ship.Movement{
		Forward:  input.NewKey(glfw.KeyW),
		Backward: input.NewKey(glfw.KeyS),

		Up:   input.NewKey(glfw.KeySpace),
		Down: input.NewKey(glfw.KeyLeftControl),

		Right: input.NewKey(glfw.KeyD),
		Left:  input.NewKey(glfw.KeyA),

		RollLeft:  input.NewKey(glfw.KeyQ),
		RollRight: input.NewKey(glfw.KeyE),
	}

	return sceneB
}

func (scene *SceneB) Render(program *program.Program) error {
	movement, roll := scene.person.Movement.UpdateMovement(program.Window, scene.camera.Dir, scene.camera.Up)

	scene.person.ApplyForce(movement)
	//scene.person.Update(deltaTime * 2)
	scene.camera.Pos.CopyVec(scene.person.Position)
	scene.camera.RotateZ(roll)

	return nil
}

func (m *SceneB) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		fmt.Println("Left mouse button pressed")
	}

	if button == glfw.MouseButtonRight && action == glfw.Press {
		fmt.Println("Right mouse button pressed")
	}

	if button == glfw.MouseButtonMiddle && action == glfw.Press {
		fmt.Println("Middle mouse button pressed")
	}
}

func (m *SceneB) MousePositionCallback(window *glfw.Window, xpos float64, ypos float64) {
	mouseX, mouseY := xpos/float64(m.Controller.ScreenWidth), ypos/float64(m.Controller.ScreenHeight)
	mouseX = mouseX*2 - 1
	mouseY = (mouseY*2 - 1)
	fmt.Println(mouseX, mouseY)

	// rotate the camera based on mouse movement
	m.camera.Rotate(mouseX, mouseY)

	window.SetCursorPos(float64(m.Controller.ScreenWidth)/2, float64(m.Controller.ScreenHeight)/2)
}

func (m *SceneB) CreateDataTexture() []uint8 {
	width, height := 1, 1
	RND := make([]float32, width*height*4)
	for i := 0; i < width*height*4; i++ {
		RND[i] = rand.Float32()
	}

	source := rand.NewSource(5)
	r := rand.New(source)

	pixels := make([]uint8, width*height*4) // 4 for RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := (y*width + x) * 4
			pixels[offset] = uint8(r.Intn(8))   // Red
			pixels[offset+1] = uint8(r.Intn(8)) // Green
			pixels[offset+2] = uint8(r.Intn(8)) // Blue
			pixels[offset+3] = 255              // Alpha
		}
	}

	return pixels
}

func (m *SceneB) Lights() *program.Light {
	return m.light
}

func (m *SceneB) Camera() *program.Camera {
	return m.camera
}
