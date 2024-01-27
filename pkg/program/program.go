package program

import (
	glprogram "remnant/pkg/gl"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	TIME_UNFOIRM_NAME       = "time\x00"
	CAMERA_POS_UNIFORM_NAME = "camera_position\x00"
	CAMERA_DIR_UNIFORM_NAME = "camera_direction\x00"
	CAMERA_UP_UNIFORM_NAME  = "camera_up\x00"
	CAMERA_FOV_UNIFORM_NAME = "camera_fov\x00"
	LIGHT_POS_UNIFORM_NAME  = "light\x00"
	RESOLUTION_UNIFORM_NAME = "resolution\x00"
	DATA_UNIFORM_NAME       = "tex\x00"
)

type Program struct {
	GLProgram *glprogram.GLProgram
	Window    *glfw.Window

	VAO        uint32
	time       int32
	cameraPos  int32
	cameraDir  int32
	cameraUp   int32
	cameraFOV  int32
	lightPos   int32
	resolution int32
	texture    int32
}

var vertices = []float32{
	// First Triangle
	-1.0, -1.0, 0.0, 0.0, 0.0, // Bottom Left
	1.0, -1.0, 0.0, 1.0, 0.0, // Bottom Right
	-1.0, 1.0, 0.0, 0.0, 1.0, // Top Left

	1.0, -1.0, 0.0, 1.0, 0.0, // Bottom Right
	1.0, 1.0, 0.0, 1.0, 1.0, // Top Right
	-1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
}

func NewProgram(windows *glfw.Window) *Program {
	program, err := glprogram.CreateGLProgram()
	if err != nil {
		panic(err)
	}

	s := &Program{
		GLProgram:  program,
		Window:     windows,
		VAO:        createTriangleVAO(vertices),
		time:       gl.GetUniformLocation(program.Handle, gl.Str(TIME_UNFOIRM_NAME)),
		cameraPos:  gl.GetUniformLocation(program.Handle, gl.Str(CAMERA_POS_UNIFORM_NAME)),
		cameraDir:  gl.GetUniformLocation(program.Handle, gl.Str(CAMERA_DIR_UNIFORM_NAME)),
		cameraUp:   gl.GetUniformLocation(program.Handle, gl.Str(CAMERA_UP_UNIFORM_NAME)),
		cameraFOV:  gl.GetUniformLocation(program.Handle, gl.Str(CAMERA_FOV_UNIFORM_NAME)),
		lightPos:   gl.GetUniformLocation(program.Handle, gl.Str(LIGHT_POS_UNIFORM_NAME)),
		resolution: gl.GetUniformLocation(program.Handle, gl.Str(RESOLUTION_UNIFORM_NAME)),
		texture:    gl.GetUniformLocation(program.Handle, gl.Str(DATA_UNIFORM_NAME)),
	}

	program.Use()
	return s
}

func (s *Program) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (s *Program) Draw() {
	gl.BindVertexArray(s.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (s *Program) SetObjectsTextureData(data []uint8) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, int32(len(data)/4), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))

	return texture
}

func (s *Program) SetClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (s *Program) SetTime(time float32) {
	gl.Uniform1f(s.time, time)
}

func (s *Program) SetCamera(camera *Camera) {
	gl.Uniform3f(s.cameraPos, float32(camera.Pos.AtVec(0)), float32(camera.Pos.AtVec(1)), float32(camera.Pos.AtVec(2)))
	gl.Uniform3f(s.cameraDir, float32(camera.Dir.AtVec(0)), float32(camera.Dir.AtVec(1)), float32(camera.Dir.AtVec(2)))
	gl.Uniform3f(s.cameraUp, float32(camera.Up.AtVec(0)), float32(camera.Up.AtVec(1)), float32(camera.Up.AtVec(2)))
	gl.Uniform1f(s.cameraFOV, camera.FOV)
}

func (s *Program) SetLight(light *Light) {
	gl.Uniform3f(s.lightPos, float32(light.Position.AtVec(0)), float32(light.Position.AtVec(1)), float32(light.Position.AtVec(2)))
}

func (s *Program) SetResolution(width, height int) {
	gl.Uniform2f(s.resolution, float32(width), float32(height))
}

func (s *Program) SetData(data uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, data)
	gl.Uniform1i(s.texture, 0)
}

func (s *Program) Delete() {
	gl.DeleteVertexArrays(1, &s.VAO)
	s.GLProgram.Delete()
}

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
