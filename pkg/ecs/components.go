package ecs

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type (
	MeshComponent struct {
	}

	RenderComponent struct {
	}

	TransformComponent struct {
		Pos mgl.Vec3
		Rot mgl.Quat
	}
)

const (
	TF_COMPID ComponentID = iota
	VEL_COMPID
	RENDER_COMPID
	MESH_COMPID
	EVENTLISTENER_COMPID
)
