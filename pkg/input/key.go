package input

import "github.com/go-gl/glfw/v3.3/glfw"

type Key struct {
	Pressed bool
	Key     glfw.Key
	window  *glfw.Window
}

func NewKey(key glfw.Key) *Key {
	return &Key{
		Pressed: false,
		Key:     key,
	}
}

func (k *Key) UpdateKeyState(window *glfw.Window) bool {
	action := window.GetKey(k.Key)

	if action == glfw.Press {
		k.Pressed = true
	}
	if action == glfw.Release {
		k.Pressed = false
	}

	return k.Pressed
}
