package ecs

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type (
	ComponentID uint32

	MeshComponent struct {
	}

	RenderComponent struct {
	}

	TransformComponent struct {
		Pos mgl.Vec3
		Rot mgl.Quat
	}

	InputReactorComponent struct {
	}
)

const (
	TF_COMPID ComponentID = iota
	RENDER_COMPID
	MESH_COMPID
	EVENTLISTENER_COMPID
)
