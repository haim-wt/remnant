package scene

import (
	"remnant/pkg/program"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Scene interface {
	Render(program *program.Program) error
	MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey)
	MousePositionCallback(window *glfw.Window, xpos float64, ypos float64)
	CreateDataTexture() []uint8
	Lights() *program.Light
	Camera() *program.Camera
}
