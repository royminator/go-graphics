package gfx

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	"strings"
)

type (
	Shader struct {
		handle   uint32
		uniforms map[string]Uniform
	}

	getGlStatus  func(uint32, uint32, *int32)
	getGlInfoLog func(uint32, int32, *int32, *uint8)
)

const ()

func CompileShader(shaderType uint32, source string) (uint32, error) {
	handle := gl.CreateShader(shaderType)
	glSrcs, freeFn := gl.Strs(source, "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrcs, nil)
	gl.CompileShader(handle)

	err := checkError(handle, gl.COMPILE_STATUS, gl.GetShaderiv,
		gl.GetShaderInfoLog, "failed to compile shader")
	if err != nil {
		return 0, err
	}

	return handle, nil
}

func CreateShaderProgram(shaders []uint32) (uint32, error) {
	prog := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(prog, shader)
	}
	gl.LinkProgram(prog)

	err := checkError(prog, gl.LINK_STATUS, gl.GetProgramiv,
		gl.GetProgramInfoLog, "failed to link program")
	if err != nil {
		return 0, err
	}

	return prog, nil
}

func checkError(handle uint32, glProp uint32, getStatusFn getGlStatus, getLogFn getGlInfoLog, msg string) error {
	var success int32
	getStatusFn(handle, glProp, &success)
	if success == gl.FALSE {
		var logLen int32
		getStatusFn(handle, gl.INFO_LOG_LENGTH, &logLen)

		log := gl.Str(strings.Repeat("\x00", int(logLen)))
		getLogFn(handle, logLen, nil, log)

		return fmt.Errorf("%s: %s", msg, gl.GoStr(log))
	}

	return nil
}

func LoadVertFragFromFile(dir string, name string) (Shader, error) {
	vsSrc, err := ioutil.ReadFile(dir + name + ".vert")
	if err != nil {
		err = fmt.Errorf("gfx: error loading vertex shader '%s': %v", name, err)
		return Shader{}, err
	}

	fsSrc, err := ioutil.ReadFile(dir + name + ".frag")
	if err != nil {
		err = fmt.Errorf("gfx: error loading fragment shader '%s': %v", name, err)
		return Shader{}, err
	}

	vs, err := CompileShader(gl.VERTEX_SHADER, string(vsSrc))
	if err != nil {
		err = fmt.Errorf("gfx: error compiling shader: %v", err)
		return Shader{}, err
	}

	fs, err := CompileShader(gl.FRAGMENT_SHADER, string(fsSrc))
	if err != nil {
		err = fmt.Errorf("gfx: error compiling shader: %v", err)
		return Shader{}, nil
	}

	shader, err := CreateShaderProgram([]uint32{vs, fs})
	if err != nil {
		err = fmt.Errorf("gfx: error linking shaders %s: %v", name, err)
		return Shader{}, nil
	}

	return Shader{handle: shader, uniforms: map[string]Uniform{}}, nil
}

func (shader *Shader) Bind() {
	gl.UseProgram(shader.handle)
	shader.bindUniforms()
}

func (shader *Shader) bindUniforms() {
	for _, uf := range shader.uniforms {
		uf.Bind(shader.handle)
	}
}

func (shader *Shader) SetUniform(name string, uf Uniform) {
	shader.uniforms[name] = uf
}
