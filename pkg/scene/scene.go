package scene

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Scene interface {
	Render(window *glfw.Window) error
	KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey)
	MousePositionCallback(window *glfw.Window, xpos float64, ypos float64)
}
