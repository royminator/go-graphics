package gfx

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type (
	Shader struct {
		Handle uint32
	}
)

func CompileShader(shaderType uint32, source string) (uint32, error) {
	handle := gl.CreateShader(shaderType)
	glSrcs, freeFn := gl.Strs(source, "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrcs, nil)
	gl.CompileShader(handle)

	var success int32
	gl.GetShaderiv(handle, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(handle, gl.INFO_LOG_LENGTH, &logLen)
		log := gl.Str(strings.Repeat("\x00", int(logLen)))
		gl.GetShaderInfoLog(handle, logLen, nil, log)
		return 0, fmt.Errorf("%s: %s", "", gl.GoStr(log))
	}

	return handle, nil
}

func CreateShaderProgram(shaders []uint32) (uint32, error) {
	prog := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(prog, shader)
	}
	gl.LinkProgram(prog)

	var success int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLen)
		log := gl.Str(strings.Repeat("\x00", int(logLen)))
		gl.GetProgramInfoLog(prog, logLen, nil, log)
		return 0, fmt.Errorf("%s: %s", "", gl.GoStr(log))
	}

	return prog, nil
}

func CreateVertexBufferG(vdata VertexData) uint32 {
	var buf uint32
	gl.GenBuffers(1, &buf)
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, vdata.NBytes, gl.Ptr(vdata.Data), gl.STATIC_DRAW)

	return buf
}

func CreateVertexArray(vBuf uint32, cBuf uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vBuf)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, cBuf)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)

	return vao
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

	return Shader{Handle: shader}, nil
}
