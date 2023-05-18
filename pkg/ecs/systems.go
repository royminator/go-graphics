package ecs

import (
	"sort"

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

	Archetype struct {
		id         ArchetypeID
		components []ComponentID
	}

	ArchetypeStatus struct {
		id       ArchetypeID
		isActive bool
	}

	ArchetypeRepo struct {
		alloc        IDallocator
		archEntities []ArchetypeStatus
		archComps    map[ArchetypeID][]ComponentID
	}

	ArchetypeID     uint32
	ComponentID     uint32
	ArchetypeMatrix [][]uint32

	Scene struct {
		entities   EntityComponents
		components ComponentStorage
		archetypes ArchetypeRepo
	}

	ComponentRepo struct {
		meshComps   []MeshComponent
		renderComps []RenderComponent
		tfComps     []TransformComponent
		velComps    []VelocityComponent
		eventComps  []EventListenerComponent
	}

    ComponentStorage interface {
        getComponent(entityID uint32, compID ComponentID) Component
        expand(n uint)
    }

    Component interface {
    }

	EntityComponents struct {
		nEntities   uint
		nComponents uint
		allocator   EntityIDallocator
	}

	IDallocator struct {
		ids  []uint32
		free []uint32
	}

	EntityIDallocator struct {
		IDallocator
		entities []AllocatorEntry
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
		Vel mgl.Vec3
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
	for i, archetype := range scene.archetypes.archEntities {
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
		archetypes: newArchetypeRepo(nEntities),
	}
}

func newArchetypeRepo(n uint) ArchetypeRepo {
	return ArchetypeRepo{
		alloc:        newIDallocator(n),
		archEntities: make([]ArchetypeStatus, n),
		archComps:    make(map[ArchetypeID][]ComponentID),
	}
}

func newEntities(nEntities uint) EntityComponents {
	return EntityComponents{
		nEntities:   nEntities,
		nComponents: MAX_COMPONENTS,
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
	scene.archetypes.append()
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
	if int(entity.index) >= len(scene.components.tfComps) {
		panic("error: couldn't add TF comp. Entity index out of bounds of component storage")
	}
	scene.components.tfComps[entity.index] = comp
	scene.archetypes.addComponent(entity.index, TF_COMPID)
}

func (scene *Scene) AddVelComp(entity EntityID, comp VelocityComponent) {
	if int(entity.index) >= len(scene.components.velComps) {
		panic("error: couldn't add Vel comp. Entity index out of bounds of component storage")
	}
	scene.components.velComps[entity.index] = comp
	scene.archetypes.addComponent(entity.index, VELOCITY_COMPID)
}

func AddComponent[T component](scene *Scene, entity EntityID, compID ComponentID, comp T) {
	comps := getComponents[T](&scene.components, compID)
	comps.addComponent(entity, comp)
	scene.archetypes.addComponent(entity.index, compID)
}

func getComponents[T component](repo *ComponentRepo, compID ComponentID) *ComponentCollection[T] {
	switch compID {
	case TF_COMPID:
		return &repo.tf2Comps
	}
}

func (repo *ArchetypeRepo) addComponent(entity uint32, comp ComponentID) {
	currArchetype := repo.archEntities[entity]
	comps := repo.archComps[currArchetype.id]
	comps = append(comps, comp)
	newArchetypeID := repo.getArchetype(comps).id
	repo.archEntities[entity].id = newArchetypeID
	repo.archEntities[entity].isActive = true
}

func (repo *ArchetypeRepo) getArchetype(comps []ComponentID) Archetype {
	sort.Slice(comps, func(i, j int) bool { return comps[i] < comps[j] })

	archetype, exists := repo.archetypeIDFromComponents(comps)
	if !exists {
		archetype = ArchetypeID(repo.alloc.allocate())
	}

	return Archetype{
		id:         archetype,
		components: comps,
	}
}

func (repo *ArchetypeRepo) archetypeIDFromComponents(comps []ComponentID) (ArchetypeID, bool) {
	for k, v := range repo.archComps {
		if sliceEquals(comps, v) {
			return k, true
		}
	}
	return 0, false
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

func (repo *ArchetypeRepo) append() {
	repo.archEntities = append(repo.archEntities, ArchetypeStatus{0, false})
}
