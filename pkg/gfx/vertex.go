package gfx

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

type VertexData struct {
	Data   interface{}
	NBytes int
}

func FromInt(v []int32) VertexData {
    return VertexData{Data: &v[0], NBytes: len(v)*4}
}

func FromFloat(v []float32) VertexData {
    return VertexData{Data: &v[0], NBytes: len(v)*4}
}

func FromVec2(v []mgl.Vec2) VertexData {
    return VertexData{Data: &v[0][0], NBytes: len(v)*2*4}
}

func FromVec3(v []mgl.Vec3) VertexData {
    return VertexData{Data: &v[0][0], NBytes: len(v)*3*4}
}

func FromVec4(v []mgl.Vec4) VertexData {
    return VertexData{Data: &v[0][0], NBytes: len(v)*4*4}
}

