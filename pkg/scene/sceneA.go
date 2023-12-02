package scenes

import (
	"fmt"
	"math/rand"
	"remnant/pkg/program"

	"github.com/go-gl/glfw/v3.3/glfw"
	"gonum.org/v1/gonum/mat"
)

var scene = program.NewScene()

func SceneA(window *glfw.Window) error {
	// Create the shader program
	program := program.NewProgram()
	defer program.Delete()

	// Load the texture data to the GPU
	data := program.SetObjectsTextureData(createDataTexture())

	// Set the clear color to black
	program.SetClearColor(0.0, 0.0, 0.0, 1.0)

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
func createDataTexture() []uint8 {
	width, height := 1, 128
	RND := make([]float32, width*height*4)
	for i := 0; i < width*height*4; i++ {
		RND[i] = rand.Float32()
	}
	rand.Seed(5)
	pixels := make([]uint8, width*height*4) // 4 for RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := (y*width + x) * 4
			pixels[offset] = uint8(rand.Intn(256))   // Red
			pixels[offset+1] = uint8(rand.Intn(256)) // Green
			pixels[offset+2] = uint8(rand.Intn(256)) // Blue
			pixels[offset+3] = 255                   // Alpha
		}
	}

	return pixels
}
func KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	right := mat.NewVecDense(3, []float64{1, 0, 0})
	up := mat.NewVecDense(3, []float64{0, 1, 0})
	forward := mat.NewVecDense(3, []float64{0, 0, 1})

	if key == glfw.KeyW && action == glfw.Press {
		scene.Camera.Pos.AddVec(scene.Camera.Pos, forward)
	}
	if key == glfw.KeyS && action == glfw.Press {
		scene.Camera.Pos.SubVec(scene.Camera.Pos, forward)
	}

	if key == glfw.KeySpace && action == glfw.Press {
		scene.Camera.Pos.AddVec(scene.Camera.Pos, up)
	}
	if key == glfw.KeyC && action == glfw.Press {
		scene.Camera.Pos.SubVec(scene.Camera.Pos, up)
	}

	if key == glfw.KeyD && action == glfw.Press {
		scene.Camera.Pos.AddVec(scene.Camera.Pos, right)
	}
	if key == glfw.KeyA && action == glfw.Press {
		scene.Camera.Pos.SubVec(scene.Camera.Pos, right)
	}

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}
