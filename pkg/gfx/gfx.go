package gfx

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	Handle uint32
}

func CompileShader(shaderType uint32, source string) (uint32, error) {
	handle := gl.CreateShader(shaderType)
	glSrcs, freeFn := gl.Strs(source, "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrcs, nil)
	gl.CompileShader(handle)

	err := checkError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "failed to like shader")
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

	err := checkError(prog, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog, "failed to like program")
	if err != nil {
		return 0, err
	}

	return prog, nil
}

type getGlStatus func(uint32, uint32, *int32)
type getGlInfoLog func(uint32, int32, *int32, *uint8)

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
