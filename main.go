package main

import (
	"fmt"
	"log"
	"remnant/internal/controller"
	"remnant/pkg/game"
	"remnant/pkg/scene"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800 * 2
	windowHeight = 450 * 2
)

// Initialization

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}
func main() {
	// Create the window
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	// OpenGL version 4.1 Core Profile
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Remnant", nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()

	// Initialize OpenGL
	err = gl.Init()
	if err != nil {
		log.Fatalln("failed to initialize OpenGL:", err)
	}

	// Run the Game
	gameController := controller.NewGameController(windowWidth, windowHeight)
	game := game.NewGame(window, gameController)
	scene := scene.NewSceneA(gameController)

	err = game.LoadScene(scene)
	if err != nil {
		log.Fatal(err)
	}

	// Terminate GLFW
	glfw.Terminate()
}
