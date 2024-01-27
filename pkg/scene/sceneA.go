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

var deltaTime = 0.0

type SceneA struct {
	*controller.Controller
	camera *program.Camera
	light  *program.Light
	ship   *ship.Ship
}

func NewSceneA(ctr *controller.Controller) *SceneA {
	sceneA := &SceneA{
		Controller: ctr,
		light:      program.NewLight(mat.NewVecDense(3, []float64{0, 1000, 1000})),
		camera:     program.NewCamera(mat.NewVecDense(3, []float64{0, 128, 64}), 90),
		ship:       ship.NewShip(mat.NewVecDense(3, []float64{0, 128, 64})),
	}

	sceneA.ship.Movement = &ship.Movement{
		Forward:  input.NewKey(glfw.KeyW),
		Backward: input.NewKey(glfw.KeyS),

		Up:   input.NewKey(glfw.KeySpace),
		Down: input.NewKey(glfw.KeyLeftControl),

		Right: input.NewKey(glfw.KeyD),
		Left:  input.NewKey(glfw.KeyA),

		RollLeft:  input.NewKey(glfw.KeyQ),
		RollRight: input.NewKey(glfw.KeyE),
	}

	return sceneA
}

func (scene *SceneA) Render(program *program.Program) error {
	movement, roll := scene.ship.Movement.UpdateMovement(program.Window, scene.camera.Dir, scene.camera.Up)

	scene.ship.ApplyForce(movement)
	scene.ship.Update(deltaTime * 2)
	scene.camera.Pos.CopyVec(scene.ship.Position)
	scene.camera.RotateZ(roll)

	return nil
}

func (m *SceneA) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
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

func (m *SceneA) MousePositionCallback(window *glfw.Window, xpos float64, ypos float64) {
	mouseX, mouseY := xpos/float64(m.Controller.ScreenWidth), ypos/float64(m.Controller.ScreenHeight)
	mouseX = mouseX*2 - 1
	mouseY = (mouseY*2 - 1)
	fmt.Println(mouseX, mouseY)

	// rotate the camera based on mouse movement
	m.camera.Rotate(mouseX, mouseY)

	window.SetCursorPos(float64(m.Controller.ScreenWidth)/2, float64(m.Controller.ScreenHeight)/2)
}

func (m *SceneA) CreateDataTexture() []uint8 {
	width, height := 1, 64
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
			pixels[offset] = uint8(r.Intn(64))   // Red
			pixels[offset+1] = uint8(r.Intn(64)) // Green
			pixels[offset+2] = uint8(r.Intn(64)) // Blue
			pixels[offset+3] = 255               // Alpha
		}
	}

	return pixels
}

func (m *SceneA) Lights() *program.Light {
	return m.light
}

func (m *SceneA) Camera() *program.Camera {
	return m.camera
}
