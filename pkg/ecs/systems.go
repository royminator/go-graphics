package ecs

import ()

type (
	EntityID        int
	ComponentMask   uint32
	ComponentID     ComponentMask
	ComponentMatrix []ComponentMask

	Scene struct {
		entities   EntityComponents
		components ComponentRepo
	}

	ComponentRepo struct {
		meshComps   []MeshComponent
		renderComps []RenderComponent
		tfComps     []TransformComponent
		velComps    []VelocityComponent
	}

	EntityComponents struct {
		nEntities   uint32
		nComponents uint32
		matrix      ComponentMatrix
	}

	MeshComponent struct {
	}

	RenderComponent struct {
	}

	TransformComponent struct {
	}

	VelocityComponent struct {
	}

	System interface {
		update(dt float32, entities []EntityID, compRepo *ComponentRepo)
		mask() ComponentMask
	}

	MovementSystem struct {
		components ComponentMask
	}
)

const (

	// If exceeds 32 different entities, change the type of ComponentMask
	TF_COMPID ComponentID = iota + 1
	VELOCITY_COMPID
	RENDER_COMPID
	MESH_COMPID
)

var (
	SYSTEMS = []System{
		MovementSystem{components: makeComponentMask(TF_COMPID, VELOCITY_COMPID)},
	}
)

func (scene *Scene) Update(dt float32) {
	for _, system := range SYSTEMS {
		entities := scene.filterEntities(system.mask())
		system.update(dt, entities, &scene.components)
	}
}

func (sys MovementSystem) update(dt float32, entities []EntityID,
	compRepo *ComponentRepo) {
}

func (sys MovementSystem) mask() ComponentMask {
	return sys.components
}

func (scene *Scene) filterEntities(mask ComponentMask) []EntityID {
	var entities []EntityID
	for i, entityMask := range scene.entities.matrix {
		if entityMask == mask {
			entities = append(entities, EntityID(i))
		}
	}
	return entities
}

func makeComponentMask(comps ...ComponentID) ComponentMask {
	var mask ComponentMask
	for _, c := range comps {
		mask = mask | ComponentMask(c)
	}
	return mask
}
