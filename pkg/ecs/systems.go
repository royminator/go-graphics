package ecs

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

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
		eventComps  []EventListenerComponent
	}

	EntityComponents struct {
		nEntities   uint
		nComponents uint
		matrix      ComponentMatrix
		allocator   EntityIDAllocator
	}

	EntityIDAllocator struct {
		IDallocator
		entities []AllocatorEntry
	}

	IDallocator struct {
		ids  []uint32
		free []uint32
	}

	MeshComponent struct {
	}

	RenderComponent struct {
	}

	TransformComponent struct {
		Pos mgl.Vec3
		Rot mgl.Quat
	}

	VelocityComponent struct {
	}

	EventListenerComponent struct {
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

	// If exceeds 32-1 different components, change the type of ComponentMask
	TF_COMPID ComponentID = iota + 1
	VELOCITY_COMPID
	RENDER_COMPID
	MESH_COMPID
	EVENTLISTENER_COMPID
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
		entities:    make([]AllocatorEntry, n),
		IDallocator: newIDallocator(n),
	}
	return alloc
}

func newIDallocator(n uint) IDallocator {
	alloc := IDallocator{
		ids:  make([]uint32, n),
		free: make([]uint32, n),
	}
	for i := range alloc.free {
		alloc.free[i] = uint32(i)
	}
	return alloc
}

func (scene *Scene) NewEntity() EntityID {
	if scene.numFreeEntities() <= 0 {
		return scene.appendEntity()
	}
	return scene.entities.allocator.allocate()
}

func (scene *Scene) appendEntity() EntityID {
	entities := &scene.entities
	alloc := &entities.allocator
	entities.nEntities++
	entities.matrix = append(entities.matrix, ComponentMask(0))
	scene.components.append()
	return alloc.append()
}

func (scene *Scene) numFreeEntities() uint {
	return uint(len(scene.entities.allocator.free))
}

func (scene *Scene) DeleteEntity(entity EntityID) {
	entities := &scene.entities
	alloc := &entities.allocator
	alloc.deallocate(entity)
}

func (alloc *EntityIDAllocator) allocate() EntityID {
	entityIndex := alloc.IDallocator.allocate()
	alloc.entities[entityIndex].version++
	alloc.entities[entityIndex].isActive = true
	return EntityID{
		index:   entityIndex,
		version: alloc.entities[entityIndex].version,
	}
}

func (alloc *IDallocator) allocate() uint32 {
	var i uint32
	i, alloc.free = alloc.free[0], alloc.free[1:]
	return i
}

func (alloc *IDallocator) deallocate(id uint32) {
	alloc.free = append(alloc.free, id)
}

func (alloc *EntityIDAllocator) append() EntityID {
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
		alloc.IDallocator.deallocate(entity.index)
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

func (comps *ComponentRepo) append() {
	comps.meshComps = append(comps.meshComps, MeshComponent{})
	comps.tfComps = append(comps.tfComps, TransformComponent{})
	comps.renderComps = append(comps.renderComps, RenderComponent{})
	comps.velComps = append(comps.velComps, VelocityComponent{})
}

func (scene *Scene) AddTfComp(entity EntityID, comp TransformComponent) {
	scene.components.tfComps[entity.index] = comp
}
