package game

import (
	"fmt"
	"image"
	"image/color"
	"remnant/internal/controller"
	"remnant/pkg/program"
	"remnant/pkg/scene"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Game struct {
	*controller.Controller
	Window *glfw.Window
}

func NewGame(window *glfw.Window) *Game {
	window.MakeContextCurrent()
	window.SetCursor(createCursor())

	w, h := window.GetFramebufferSize()
	ctrl := controller.NewGameController(w, h)

	return &Game{
		Window:     window,
		Controller: ctrl,
	}
}

func (g *Game) Load(scene scene.Scene) error {
	g.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	g.Window.SetMouseButtonCallback(scene.MouseButtonCallback)
	g.Window.SetCursorPosCallback(scene.MousePositionCallback)
	g.Window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			g.Window.SetShouldClose(true)
		}
		g.IsReady = true
	})

	return nil
}

func (g *Game) Run(window *glfw.Window, scene *scene.SceneB) error {
	// Create the shader program
	program := program.NewProgram(window)
	defer program.Delete()

	// Load the texture data to the GPU
	tx := scene.CreateDataTexture()
	data := program.SetObjectsTextureData(tx)

	// Set the clear color to black
	program.SetClearColor(0.0, 0.0, 0.0, 1.0)

	// initialize mouse position to middle of screen
	window.SetCursorPos(float64(g.ScreenWidth)/2, float64(g.ScreenHeight)/2)

	deltaTime := 0.0
	for !g.Window.ShouldClose() {

		// Statistics
		seconds := 0.0
		fps := 0.0

		program.Clear()

		scene.Render(program)

		// Update the shader uniforms
		program.SetTime(float32(seconds))
		program.SetResolution(g.ScreenWidth, g.ScreenHeight)
		program.SetData(data)
		program.SetLight(scene.Lights())
		program.SetCamera(scene.Camera())

		// Draw
		program.Draw()

		// Swap the buffers
		window.SwapBuffers()
		glfw.PollEvents()

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
