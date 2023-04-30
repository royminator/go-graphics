package ecs

import ()

type (
	EntityID struct {
		index   uint32
		version uint32
	}

	AllocatorEntry struct {
		isActive bool
		version  uint32
	}

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
		nEntities   uint
		nComponents uint
		matrix      ComponentMatrix
		allocator   EntityIDAllocator
	}

	EntityIDAllocator struct {
		entities []AllocatorEntry
		free     []uint32
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
		update(dt float32, entities []uint32, compRepo *ComponentRepo)
		mask() ComponentMask
	}

	MovementSystem struct {
		components ComponentMask
	}
)

const (
	MAX_COMPONENTS uint = 32

	// If exceeds 32 different components, change the type of ComponentMask
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

func (sys MovementSystem) update(dt float32, entities []uint32,
	compRepo *ComponentRepo) {
}

func (sys MovementSystem) mask() ComponentMask {
	return sys.components
}

func (scene *Scene) filterEntities(mask ComponentMask) []uint32 {
	var entities []uint32
	for i, entityMask := range scene.entities.matrix {
		if entityMask == mask {
			entities = append(entities, uint32(i))
		}
	}
	return entities
}

func makeComponentMask(comps ...ComponentID) ComponentMask {
	var mask ComponentMask
	for _, c := range comps {
		mask |= ComponentMask(c)
	}
	return mask
}

func NewScene(nEntities uint) *Scene {
	return &Scene{
		entities:   newEntities(nEntities),
		components: newComponents(nEntities),
	}
}

func newEntities(nEntities uint) EntityComponents {
	return EntityComponents{
		nEntities:   nEntities,
		nComponents: MAX_COMPONENTS,
		matrix:      newComponentMatrix(nEntities, MAX_COMPONENTS),
		allocator:   newEntityAllocator(nEntities),
	}
}

func newComponents(nEnts uint) ComponentRepo {
	return ComponentRepo{
		meshComps:   make([]MeshComponent, nEnts),
		renderComps: make([]RenderComponent, nEnts),
		tfComps:     make([]TransformComponent, nEnts),
		velComps:    make([]VelocityComponent, nEnts),
	}
}

func newComponentMatrix(nEnts uint, nComps uint) ComponentMatrix {
	var mat ComponentMatrix = make([]ComponentMask, nEnts)
	for i := range mat {
		mat[i] = ComponentMask(0)
	}
	return mat
}

func newEntityAllocator(n uint) EntityIDAllocator {
	alloc := EntityIDAllocator{
		entities: make([]AllocatorEntry, n),
		free:     make([]uint32, n),
	}
	for i := range alloc.free {
		alloc.free[i] = uint32(i)
	}
	return alloc
}

func (scene *Scene) NewEntity() EntityID {
	return scene.entities.allocator.allocate()
}

func (alloc *EntityIDAllocator) allocate() EntityID {
	if len(alloc.free) > 0 {
		var entityIndex uint32
		entityIndex, alloc.free = alloc.free[0], alloc.free[1:]
		alloc.entities[entityIndex].version++
		alloc.entities[entityIndex].isActive = true
		return EntityID{
			index:   entityIndex,
			version: alloc.entities[entityIndex].version,
		}
	}
	entry := AllocatorEntry{
		isActive: true,
		version:  0,
	}
	alloc.entities = append(alloc.entities, entry)
	return EntityID{
		index:   uint32(len(alloc.entities) - 1),
		version: 0,
	}
}

func (alloc *EntityIDAllocator) deallocate(entity EntityID) {
	if alloc.isActive(entity) {
		alloc.entities[entity.index].isActive = false
		alloc.free = append(alloc.free, entity.index)
	}
}

func (alloc *EntityIDAllocator) isActive(entity EntityID) bool {
	return entity.index < uint32(len(alloc.entities)) &&
		alloc.entities[entity.index].version == entity.version &&
		alloc.entities[entity.index].isActive
}

func contains[T comparable](xs []T, x T) bool {
	for _, val := range xs {
		if val == x {
			return true
		}
	}
	return false
}
