package glprogram

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	handle uint32
}

type GLProgram struct {
	Pointer uint32
	shaders []*Shader
}

func (shader *Shader) Delete() {
	gl.DeleteShader(shader.handle)
}

func (prog *GLProgram) Delete() {
	for _, shader := range prog.shaders {
		shader.Delete()
	}
	gl.DeleteProgram(prog.Pointer)
}

func (prog *GLProgram) Attach(shaders ...*Shader) {
	for _, shader := range shaders {
		gl.AttachShader(prog.Pointer, shader.handle)
		prog.shaders = append(prog.shaders, shader)
	}
}

func (prog *GLProgram) Use() uint32 {
	gl.UseProgram(prog.Pointer)
	return prog.Pointer
}

func (prog *GLProgram) Link() error {
	gl.LinkProgram(prog.Pointer)
	return getGlError(prog.Pointer, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
		"PROGRAM::LINKING_FAILURE")
}

func NewGLProgram(shaders ...*Shader) (*GLProgram, error) {
	prog := &GLProgram{Pointer: gl.CreateProgram()}
	prog.Attach(shaders...)

	if err := prog.Link(); err != nil {
		return nil, err
	}

	return prog, nil
}

func NewShader(src string, sType uint32) (*Shader, error) {

	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(src + "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)
	err := getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::")
	if err != nil {
		return nil, err
	}
	return &Shader{handle: handle}, nil
}

func NewShaderFromFile(file string, sType uint32) (*Shader, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(string(src) + "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)
	err = getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::"+file)
	if err != nil {
		return nil, err
	}
	return &Shader{handle: handle}, nil
}

func CreateGLProgram() (*GLProgram, error) {
	vertShader, err := NewShaderFromFile("shaders/vertex.glsl", gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragShader, err := NewShaderFromFile("shaders/fragment.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	shaderProgram, err := NewGLProgram(vertShader, fragShader)
	if err != nil {
		return nil, err
	}

	return shaderProgram, nil
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, log)

		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
	}

	return nil
}
