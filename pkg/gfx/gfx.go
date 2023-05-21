package gfx

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func CreateVertexBufferG(vdata VertexData) uint32 {
	var buf uint32
	gl.GenBuffers(1, &buf)
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, vdata.NumBytes, gl.Ptr(vdata.Data), gl.STATIC_DRAW)

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
