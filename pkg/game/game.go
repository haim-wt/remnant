package game

import (
	"image"
	"image/color"
	"remnant/internal/controller"
	"remnant/pkg/scene"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Game struct {
	*controller.GameController
	Window *glfw.Window
}

func NewGame(window *glfw.Window, ctr *controller.GameController) *Game {
	window.MakeContextCurrent()
	window.SetCursor(createCursor())

	return &Game{
		Window:         window,
		GameController: ctr,
	}
}

func (g *Game) LoadScene(scene scene.Scene) error {
	g.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	g.Window.SetMouseButtonCallback(scene.MouseButtonCallback)
	g.Window.SetCursorPosCallback(scene.MousePositionCallback)
	g.Window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			g.Window.SetShouldClose(true)
		}
		g.IsReady = true
	})
	return scene.Render(g.Window)
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
