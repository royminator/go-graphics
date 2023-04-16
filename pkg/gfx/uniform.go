package gfx

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type (
	Uniform interface {
		Bind(prog uint32)
	}

	Uniform4f struct {
		Name string
		Val  mgl.Mat4
	}
)

func (uf Uniform4f) Bind(prog uint32) {
	loc := gl.GetUniformLocation(prog, gl.Str(uf.Name+"\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &uf.Val[0])
}
