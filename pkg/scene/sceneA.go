package scene

import (
	"fmt"
	"math/rand"
	"remnant/internal/controller"
	"remnant/pkg/program"

	"github.com/go-gl/glfw/v3.3/glfw"
	"gonum.org/v1/gonum/mat"
)

type SceneA struct {
	*controller.GameController
	Camera *program.Camera
	Light  *program.Light
}

func NewSceneA(ctr *controller.GameController) *SceneA {
	return &SceneA{
		GameController: ctr,
		Light: &program.Light{
			Position: mat.NewVecDense(3, []float64{0, 64, -64}),
		},
		Camera: &program.Camera{
			Pos: mat.NewVecDense(3, []float64{-32, 0, -32}),
			Dir: mat.NewVecDense(3, []float64{0, 0, 1}),
			Up:  mat.NewVecDense(3, []float64{0, 1, 0}),
			FOV: float32(90),
		},
	}
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
	timeElapsed := 0.0
	fps := 0.0

	w, h := window.GetSize()
	for !window.ShouldClose() {
		// CPU Events
		glfw.PollEvents()

		// Clear the screen
		program.Clear()

		// Update the shader uniforms
		program.SetTime(float32(timeElapsed))
		program.SetResolution(w, h)
		program.SetData(data)
		program.SetLight(scene.Light)
		program.SetCamera(scene.Camera)

		// Draw
		program.Draw()

		// Swap the buffers
		window.SwapBuffers()

		fps += 1.0
		if glfw.GetTime() >= 1.0 {
			window.SetTitle(fmt.Sprintf("FPS: %.2f", fps))
			glfw.SetTime(0.0)
			fps = 0.0
		}
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

	// rotate the camera based on mouse movement
	m.Camera.Rotate(mouseX, mouseY)

	window.SetCursorPos(float64(m.GameController.ScreenWidth)/2, float64(m.GameController.ScreenHeight)/2)
}

func (m *SceneA) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	right := mat.NewVecDense(3, []float64{1, 0, 0})
	up := mat.NewVecDense(3, []float64{0, 1, 0})
	forward := mat.NewVecDense(3, []float64{0, 0, 1})

	if key == glfw.KeyW && action == glfw.Press {
		m.Camera.Pos.AddVec(m.Camera.Pos, forward)
	}
	if key == glfw.KeyS && action == glfw.Press {
		m.Camera.Pos.SubVec(m.Camera.Pos, forward)
	}

	if key == glfw.KeySpace && action == glfw.Press {
		m.Camera.Pos.AddVec(m.Camera.Pos, up)
	}
	if key == glfw.KeyC && action == glfw.Press {
		m.Camera.Pos.SubVec(m.Camera.Pos, up)
	}

	if key == glfw.KeyD && action == glfw.Press {
		m.Camera.Pos.AddVec(m.Camera.Pos, right)
	}
	if key == glfw.KeyA && action == glfw.Press {
		m.Camera.Pos.SubVec(m.Camera.Pos, right)
	}

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func createDataTexture() []uint8 {
	width, height := 1, 128
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
			pixels[offset] = uint8(r.Intn(256))   // Red
			pixels[offset+1] = uint8(r.Intn(256)) // Green
			pixels[offset+2] = uint8(r.Intn(256)) // Blue
			pixels[offset+3] = 255                // Alpha
		}
	}

	return pixels
}
