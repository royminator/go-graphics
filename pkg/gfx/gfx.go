package gfx

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

func CreateShader(shaderType uint32, source string) (uint32, error) {
    handle := gl.CreateShader(shaderType)
    glSrcs, freeFn := gl.Strs(source, "\x00")
    defer freeFn()
    gl.ShaderSource(handle, 1, glSrcs, nil)
    gl.CompileShader(handle)

    var status int32
    gl.GetShaderiv(handle, gl.COMPILE_STATUS, &status)
    fmt.Printf("Compiled shader of type %d with result %X\n", shaderType, status)

    return handle, nil
}

func CreateShaderProgram(shaders []uint32) uint32 {
    prog := gl.CreateProgram()

    for _, shader := range shaders {
        gl.AttachShader(prog, shader)
    }

    gl.LinkProgram(prog)
    var status int32
    gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
    fmt.Printf("Linked shader program %d with result %X\n", prog, status)

    return prog
}

func CreateVertexBuffer(data []mgl.Vec3) uint32 {
	var buf uint32
    ptr := gl.Ptr(&data[0][0])
	gl.GenBuffers(1, &buf)
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4*len(data[0]), ptr, gl.STATIC_DRAW)

    return buf
}

func CreateColorBuffer(data []mgl.Vec4) uint32 {
	var buf uint32
	gl.GenBuffers(1, &buf)
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4*len(data[0]), gl.Ptr(&data[0][0]), gl.STATIC_DRAW)

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
