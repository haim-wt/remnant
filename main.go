package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const windowWidth = 800
const windowHeight = 600

var res [2]float32
var lp [3]float32
var ro [3]float32

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
	res = [2]float32{windowWidth, windowHeight}
	lp = [3]float32{-0.1, 1.0, 1.2}
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to inifitialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "basic shaders", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(keyCallback)

	err = programLoop(window)
	if err != nil {
		log.Fatal(err)
	}
}

/*
 * Creates the Vertex Array Object for a triangle.
 */
func createTriangleVAO(vertices []float32) uint32 {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	// Texture coordinate attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	return vao
}

func createTexture() uint32 {
	width, height := 1, 32
	RND := make([]float32, width*height*4)
	for i := 0; i < width*height*4; i++ {
		RND[i] = rand.Float32()
	}
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

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	return texture
}

func createProgram(vertexShaderSource, fragmentShaderSource string) (*Program, error) {
	vertShader, err := NewShaderFromFile("vertex.glsl", gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragShader, err := NewShaderFromFile("fragment.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	shaderProgram, err := NewProgram(vertShader, fragShader)
	if err != nil {
		return nil, err
	}

	return shaderProgram, nil
}
func programLoop(window *glfw.Window) error {
	// the linked shader program determines how the data will be rendered

	program, err := createProgram("vertex.glsl", "fragment.glsl")
	if err != nil {
		return err
	}
	defer program.Delete()

	vertices := []float32{
		// First Triangle
		-1.0, -1.0, 0.0, 0.0, 0.0, // Bottom Left
		1.0, -1.0, 0.0, 1.0, 0.0, // Bottom Right
		-1.0, 1.0, 0.0, 0.0, 1.0, // Top Left

		1.0, -1.0, 0.0, 1.0, 0.0, // Bottom Right
		1.0, 1.0, 0.0, 1.0, 1.0, // Top Right
		-1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
	}

	c := 0.0
	VAO := createTriangleVAO(vertices)

	progPoiner := program.Use()

	timeLocation := gl.GetUniformLocation(progPoiner, gl.Str("time\x00"))
	resolutionLocation := gl.GetUniformLocation(progPoiner, gl.Str("resolution\x00"))
	light_positionLocation := gl.GetUniformLocation(progPoiner, gl.Str("light_position\x00"))
	origin_position := gl.GetUniformLocation(progPoiner, gl.Str("ray_origin\x00"))
	texLocation := gl.GetUniformLocation(progPoiner, gl.Str("tex\x00"))

	gl.ActiveTexture(gl.TEXTURE0)
	texture := createTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)

	for !window.ShouldClose() {
		startTime := time.Now()
		// poll events and call their registered callbacks
		glfw.PollEvents()

		// perform rendering
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// draw loop
		c -= 0.2
		ro = [3]float32{0, 0, -float32(c)}

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 6) // 6 vertices for two triangles
		gl.BindVertexArray(0)

		gl.Uniform1f(timeLocation, float32(c))
		gl.Uniform2fv(resolutionLocation, 1, &res[0])
		gl.Uniform4fv(light_positionLocation, 1, &lp[0])
		gl.Uniform3fv(origin_position, 1, &ro[0])
		gl.Uniform1i(texLocation, 0)

		window.SwapBuffers()

		frameTime := time.Since(startTime).Milliseconds()
		log.Default().Println("Frame time:", frameTime, "ms")
	}

	return nil
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}
