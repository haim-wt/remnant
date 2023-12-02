package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	scenes "remnant/pkg/scene"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 450
)

// Initialization

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
	err := scenes.SceneA(window)
	if err != nil {
		log.Fatal(err)
	}

	// Terminate GLFW
	glfw.Terminate()
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
	window.SetInputMode(glfw.StickyKeysMode, 1)
	window.SetInputMode(glfw.StickyMouseButtonsMode, 1)
	window.SetKeyCallback(scenes.KeyCallback)
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
