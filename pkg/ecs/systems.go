package ecs

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"sort"
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

	Archetype struct {
		id         ArchetypeID
		components []ComponentID
	}

	ArchetypeStatus struct {
		id       ArchetypeID
		isActive bool
	}

	ArchetypeID     uint32
	ComponentID     uint32
	ArchetypeMatrix [][]uint32

	Scene struct {
		entities   EntityComponents
		components ComponentRepo
		archetypes []Archetype
	}

	ComponentRepo struct {
		meshComps   []MeshComponent
		renderComps []RenderComponent
		tfComps     []TransformComponent
		velComps    []VelocityComponent
		eventComps  []EventListenerComponent
	}

	EntityComponents struct {
		nEntities    uint
		nComponents  uint
		archEntities []ArchetypeStatus
		allocator    EntityIDallocator
	}

	IDallocator struct {
		ids  []uint32
		free []uint32
	}

	EntityIDallocator struct {
		IDallocator
		entities []AllocatorEntry
	}

	ArchetypeIDallocator struct {
		IDallocator
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
		archetypeID() ArchetypeID
	}

	MovementSystem struct {
		archetype ArchetypeID
	}
)

const (
	MAX_COMPONENTS uint = 32
	MAX_ARCHETYPES uint = 300

	TF_COMPID ComponentID = iota + 1
	VELOCITY_COMPID
	RENDER_COMPID
	MESH_COMPID
	EVENTLISTENER_COMPID
)

var (
	SYSTEMS = []System{}
)

func (scene *Scene) Update(dt float32) {
	for _, system := range SYSTEMS {
		entities := scene.filterEntities(system.archetypeID())
		system.update(dt, entities, &scene.components)
	}
}

func (sys MovementSystem) update(dt float32, entities []uint32,
	compRepo *ComponentRepo) {
}

func (sys MovementSystem) archetypeID() ArchetypeID {
	return sys.archetype
}

func (scene *Scene) filterEntities(arch ArchetypeID) []uint32 {
	var entities []uint32
	for i, archetype := range scene.entities.archEntities {
		if archetype.id == arch && archetype.isActive {
			entities = append(entities, uint32(i))
		}
	}
	return entities
}

func NewScene(nEntities uint) *Scene {
	return &Scene{
		entities:   newEntities(nEntities),
		components: newComponents(nEntities),
	}
}

func newEntities(nEntities uint) EntityComponents {
	return EntityComponents{
		nEntities:    nEntities,
		nComponents:  MAX_COMPONENTS,
		archEntities: make([]ArchetypeStatus, nEntities),
		allocator:    newEntityAllocator(nEntities),
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

func newEntityAllocator(n uint) EntityIDallocator {
	alloc := EntityIDallocator{
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
	scene.components.append()
	entities.archEntities = append(entities.archEntities, ArchetypeStatus{0, false})
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

func (alloc *EntityIDallocator) allocate() EntityID {
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

func (alloc *EntityIDallocator) append() EntityID {
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

func (alloc *EntityIDallocator) deallocate(entity EntityID) {
	if alloc.isActive(entity) {
		alloc.entities[entity.index].isActive = false
		alloc.IDallocator.deallocate(entity.index)
	}
}

func (alloc *EntityIDallocator) isActive(entity EntityID) bool {
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

func (scene *Scene) newArchetype(alloc *IDallocator, comps []ComponentID) Archetype {
	sort.Slice(comps, func(i, j int) bool { return comps[i] < comps[j] })

	archetype, exists := scene.archetypeIDFromComponents(comps)
	if !exists {
		archetype = ArchetypeID(alloc.allocate())
	}

	return Archetype{
		id:         archetype,
		components: comps,
	}
}

func (scene *Scene) archetypeIDFromComponents(comps []ComponentID) (ArchetypeID, bool) {
	for _, arch := range scene.archetypes {
		if arch.hasComponents(comps) {
			return arch.id, true
		}
	}
	return 0, false
}

func (arch *Archetype) hasComponents(comps []ComponentID) bool {
	return sliceEquals(arch.components, comps)
}

func sliceEquals[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
