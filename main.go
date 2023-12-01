package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"gonum.org/v1/gonum/mat"
)

const windowWidth = 800
const windowHeight = 600

var camera = &Camera{
	Pos: mat.NewVecDense(3, []float64{-32, 0, -32}),
	Dir: mat.NewVecDense(3, []float64{0, 0, 1}),
	Up:  mat.NewVecDense(3, []float64{0, 1, 0}),
	FOV: float32(60),
}

var light = &Light{
	Position: mat.NewVecDense(3, []float64{0, 64, -64}),
}

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	// Create the window
	window := createGlfwWindow()

	// Initialize OpenGL
	initializeOPenGL()

	// Run the program loop
	err := programLoop(window)
	if err != nil {
		log.Fatal(err)
	}

	// Terminate GLFW
	glfw.Terminate()
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

func programLoop(window *glfw.Window) error {
	// Create the shader program
	program := NewProgram()
	defer program.Delete()

	// Load the texture data to the GPU
	data := program.SetObjectsTextureData(createDataTexture())

	// Set the clear color to black
	program.SetClearColor(0.0, 0.0, 0.0, 1.0)

	// Statistics
	timeElapsed := 0.0
	fps := 0.0

	for !window.ShouldClose() {
		// CPU Events
		glfw.PollEvents()

		// Clear the screen
		program.Clear()

		// Update the shader uniforms
		program.SetTime(float32(timeElapsed))
		program.SetResolution(windowWidth, windowHeight)
		program.SetData(data)
		program.SetLight(light)
		program.SetCamera(camera)

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

func initializeOPenGL() {
	// Initialize OpenGL
	err := gl.Init()
	if err != nil {
		log.Fatalln("failed to initialize OpenGL:", err)
	}
}
func createGlfwWindow() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	// OpenGL version 4.1 Core Profile
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "basic shaders", nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()
	window.SetInputMode(glfw.StickyMouseButtonsMode, 1)
	window.SetKeyCallback(keyCallback)
	window.SetCursor(createCursor())

	return window
}
func createCursor() *glfw.Cursor {
	cursorImage := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for x := 0; x < 8; x++ {
		cursorImage.Set(x, 0, color.RGBA{255, 255, 255, 255})
		cursorImage.Set(x, 1, color.RGBA{255, 255, 255, 255})
		cursorImage.Set(x, 2, color.RGBA{255, 255, 255, 255})
		cursorImage.Set(x, 3, color.RGBA{255, 255, 255, 255})
		cursorImage.Set(x, 4, color.RGBA{255, 255, 255, 255})
		cursorImage.Set(x, 5, color.RGBA{255, 255, 255, 255})
		cursorImage.Set(x, 6, color.RGBA{255, 255, 255, 255})
		cursorImage.Set(x, 7, color.RGBA{255, 255, 255, 255})
	}
	cursor := glfw.CreateCursor(cursorImage, 4, 4)
	return cursor
}
func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	window.SetInputMode(glfw.StickyKeysMode, 1)

	// if key == glfw.KeyW && action == glfw.Press {
	// 	cameraPos[2] -= 1
	// }
	// if key == glfw.KeyS && action == glfw.Press {
	// 	cameraPos[2] += 1
	// }
	// if key == glfw.KeyA && action == glfw.Press {
	// 	cameraPos[0] -= 0.1
	// }
	// if key == glfw.KeyD && action == glfw.Press {
	// 	cameraPos[0] += 0.1
	// }
	// if key == glfw.KeyEscape && action == glfw.Press {
	// 	window.SetShouldClose(true)
	// }
}
