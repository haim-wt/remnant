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
	*controller.GameController
	Camera *program.Camera
	Light  *program.Light
	Ship   *ship.Ship
}

func NewSceneA(ctr *controller.GameController) *SceneA {
	sceneA := &SceneA{
		GameController: ctr,
		Light:          program.NewLight(mat.NewVecDense(3, []float64{500, 1000, 300})),
		Camera:         program.NewCamera(mat.NewVecDense(3, []float64{0, 64, -64}), 90),
		Ship:           ship.NewShip(mat.NewVecDense(3, []float64{0, 64, -64})),
	}

	sceneA.Ship.Movement = &ship.Movement{
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

func (scene *SceneA) Render(window *glfw.Window) error {
	// Create the shader program
	program := program.NewProgram()
	defer program.Delete()

	// Load the texture data to the GPU
	data := program.SetObjectsTextureData(createDataTexture())

	// Set the clear color to black
	program.SetClearColor(0.0, 0.0, 0.0, 1.0)

	// initialize mouse position to middle of screen
	window.SetCursorPos(float64(scene.GameController.ScreenWidth)/2, float64(scene.GameController.ScreenHeight)/2)

	// Statistics
	seconds := 0.0
	fps := 0.0

	w, h := window.GetSize()
	for !window.ShouldClose() {

		// CPU Events
		glfw.PollEvents()

		movement, roll := scene.Ship.Movement.UpdateMovement(window, scene.Camera.Dir, scene.Camera.Up)
		scene.Ship.ApplyForce(movement)

		scene.Ship.Update(deltaTime * 9.8)
		scene.Camera.Pos.CopyVec(scene.Ship.Position)

		scene.Camera.RotateZ(roll)
		// Clear the screen
		program.Clear()

		// Update the shader uniforms
		program.SetTime(float32(seconds))
		program.SetResolution(w, h)
		program.SetData(data)
		program.SetLight(scene.Light)
		program.SetCamera(scene.Camera)

		// Draw
		program.Draw()

		// Swap the buffers
		window.SwapBuffers()

		fps += 1.0
		deltaTime = glfw.GetTime()
		seconds += deltaTime
		if seconds >= 1.0 {
			window.SetTitle(fmt.Sprintf("FPS: %.2f", fps))
			seconds = 0
			fps = 0.0
		}
		glfw.SetTime(0.0)
	}

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
	mouseX, mouseY := xpos/float64(m.GameController.ScreenWidth), ypos/float64(m.GameController.ScreenHeight)
	mouseX = mouseX*2 - 1
	mouseY = (mouseY*2 - 1)
	fmt.Println(mouseX, mouseY)

	// rotate the camera based on mouse movement
	m.Camera.Rotate(mouseX, mouseY)

	window.SetCursorPos(float64(m.GameController.ScreenWidth)/2, float64(m.GameController.ScreenHeight)/2)
}

func createDataTexture() []uint8 {
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
