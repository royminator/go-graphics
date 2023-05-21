package ecs

import (
	"go-graphics/pkg/util"
	"sort"
)

type (
	EntityID    uint32
	ArchetypeID uint32
	ComponentID uint32

	Entity struct {
		index   EntityID
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

	Scene struct {
		entities   EntityComponents
		components ComponentRepo
		archetypes ArchetypeRepo
	}

	ComponentRepo struct {
		meshComps   []MeshComponent
		renderComps []RenderComponent
		tfComps     []TransformComponent
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

func (scene *Scene) NewEntity() Entity {
	if scene.numFreeEntities() <= 0 {
		return scene.appendEntity()
	}
	return scene.entities.allocator.allocate()
}

func (scene *Scene) appendEntity() Entity {
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

func (scene *Scene) DeleteEntity(entity Entity) {
	entities := &scene.entities
	alloc := &entities.allocator
	alloc.deallocate(entity)
}

func (alloc *EntityIDallocator) allocate() Entity {
	entityIndex := alloc.IDallocator.allocate()
	alloc.entities[entityIndex].version++
	alloc.entities[entityIndex].isActive = true
	return Entity{
		index:   EntityID(entityIndex),
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

func (alloc *EntityIDallocator) append() Entity {
	entry := AllocatorEntry{
		isActive: true,
		version:  0,
	}
	alloc.entities = append(alloc.entities, entry)
	return Entity{
		index:   EntityID(len(alloc.entities) - 1),
		version: 0,
	}
}

func (alloc *EntityIDallocator) deallocate(entity Entity) {
	if alloc.isActive(entity) {
		alloc.entities[entity.index].isActive = false
		alloc.IDallocator.deallocate(uint32(entity.index))
	}
}

func (alloc *EntityIDallocator) isActive(entity Entity) bool {
	return int(entity.index) < len(alloc.entities) &&
		alloc.entities[entity.index].version == entity.version &&
		alloc.entities[entity.index].isActive
}

func (comps *ComponentRepo) append() {
	comps.meshComps = append(comps.meshComps, MeshComponent{})
	comps.tfComps = append(comps.tfComps, TransformComponent{})
	comps.renderComps = append(comps.renderComps, RenderComponent{})
}

func (scene *Scene) AddTfComp(entity Entity, comp TransformComponent) {
	comps := scene.components.tfComps
	validateAndAddComponent(scene, comps, entity.index, comp, TF_COMPID)
}

func (scene *Scene) AddMeshComp(entity Entity, comp MeshComponent) {
	comps := scene.components.meshComps
	validateAndAddComponent(scene, comps, entity.index, comp, TF_COMPID)
}

func validateAndAddComponent[T any](s *Scene, comps []T, entity EntityID, comp T, compID ComponentID) {
	if int(entity) >= len(comps) {
		panic("error: couldn't add component. Entity index out of bounds of component storage")
	}
	comps[entity] = comp
	s.archetypes.addComponent(entity, compID)
}

func (repo *ArchetypeRepo) addComponent(entity EntityID, comp ComponentID) {
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
		if util.SliceEquals(comps, v) {
			return k, true
		}
	}
	return 0, false
}

func (repo *ArchetypeRepo) append() {
	repo.archEntities = append(repo.archEntities, ArchetypeStatus{0, false})
}