package ecs

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"go-graphics/pkg/util"
	"testing"
)

func TestEntityComponents_newEntityComponents_ShouldAddEntities(t *testing.T) {
	// Arrange
	maxEntities := uint(3)

	// Act
	entities := newEntities(maxEntities)

	// Assert
	if entities.nEntities != maxEntities {
		t.Errorf("expected %d number of entities, was %d", maxEntities, entities.nEntities)
	}
	lenAllocEntities := len(entities.allocator.entities)
	lenAllocFree := len(entities.allocator.free)
	if lenAllocEntities != lenAllocFree || uint(lenAllocEntities) != maxEntities {
		t.Errorf("number of allocator entities was not equal to number of free indices, was no equal to %d", maxEntities)
	}
}

func TestScene_NewScene_ShouldCreateEntities(t *testing.T) {
	maxEntities := uint(3)
	scene := NewScene(maxEntities)
	if len(scene.components.meshComps) != int(maxEntities) {
		t.Errorf("expected %d number of mesh components, was %d", maxEntities, len(scene.components.meshComps))
	}
	if scene.entities.nEntities != maxEntities {
		t.Errorf("expected %d number of entities, was %d", maxEntities, scene.entities.nEntities)
	}
	lenAllocEntities := len(scene.entities.allocator.entities)
	lenAllocFree := len(scene.entities.allocator.free)
	if lenAllocEntities != lenAllocFree || uint(lenAllocEntities) != maxEntities {
		t.Errorf("number of allocator entities was not equal to number of free indices, was no equal to %d", maxEntities)
	}
}

func TestScene_NewEntity_ShouldBe0(t *testing.T) {
	maxEntities := uint(3)
	scene := NewScene(maxEntities)
	entity := scene.NewEntity()
	actual := Entity{EntityID(0), 1}
	if entity != actual {
		t.Errorf("expected allocated entity to be %v, was %v", entity, actual)
	}
}

func TestScene_NewEntity_WhenEntitiesCapacityReached_ShouldAppend(t *testing.T) {
	// Arrange
	maxEntities := uint(3)
	scene := NewScene(maxEntities)

	// Act
	for i := 0; i < int(maxEntities); i++ {
		scene.NewEntity()
	}
	scene.NewEntity()

	// Assert
	nEnts := maxEntities + 1
	lenEnts := scene.entities.nEntities
	lenAllocEnts := len(scene.entities.allocator.entities)
	lenAllocFree := len(scene.entities.allocator.free)
	lenRenderComps := len(scene.components.renderComps)

	if lenEnts != nEnts {
		t.Errorf("expected n entities to be %d", nEnts)
	}
	if uint(lenAllocEnts) != nEnts {
		t.Errorf("expected n allocated entities to be %d", nEnts)
	}
	if uint(lenAllocFree) != 0 {
		t.Error("expected n free indices to be 0")
	}
	if uint(lenRenderComps) != nEnts {
		t.Errorf("expected n render components to be %d", nEnts)
	}
}

func TestScene_NewEntity_WhenEntitiesCapacityReachedThenDeallocate_NumberOfEntitiesShouldEqualCapacity(t *testing.T) {
	// Arrange
	maxEntities := uint(3)
	scene := NewScene(maxEntities)

	// Act
	for i := 0; i < int(maxEntities); i++ {
		scene.NewEntity()
	}

	entity := scene.NewEntity()
	scene.DeleteEntity(entity)

	// Assert
	expEnts := maxEntities + 1
	lenEnts := scene.entities.nEntities
	lenAllocEnts := len(scene.entities.allocator.entities)
	lenAllocFree := len(scene.entities.allocator.free)
	lenRenderComps := len(scene.components.renderComps)

	if lenEnts != expEnts {
		t.Errorf("expected n entities to be %d, was %d", expEnts, lenEnts)
	}
	if uint(lenAllocEnts) != expEnts {
		t.Errorf("expected n allocated entities to be %d, was %d", expEnts, lenAllocEnts)
	}
	if uint(lenAllocFree) != 1 {
		t.Errorf("expected n free indices to be 1, was %d", lenAllocFree)
	}
	if uint(lenRenderComps) != expEnts {
		t.Errorf("expected n render components to be %d, was %d", expEnts, lenRenderComps)
	}
}

func TestScene_AddComponent(t *testing.T) {
	// Arrange
	comp := TransformComponent{mgl.Vec3{1, 2, 3}, mgl.QuatIdent()}
	scene := NewScene(200)
	entity := scene.NewEntity()

	// Act
	scene.AddTfComp(entity, comp)

	// Assert
	expected := scene.components.tfComps[entity.index]
	if comp != expected {
		t.Errorf("expected component to be %v, was %v", comp, expected)
	}
}

func TestScene_AddComponent_ShouldAddArchetypeToEntity(t *testing.T) {
	// Arrange
	comp := TransformComponent{mgl.Vec3{1, 2, 3}, mgl.QuatIdent()}
	scene := NewScene(200)
	entity := scene.NewEntity()

	// Act
	scene.AddTfComp(entity, comp)

	// Assert
	actual := scene.archetypes.archEntities[entity.index]
	expected := Archetype{0, []ComponentID{TF_COMPID}}
	status := scene.archetypes.archEntities[entity.index]
	if actual.id != expected.id {
		t.Errorf("expected archetype.id to be %d, was %d", expected.id, actual.id)
	}
	if status.id != expected.id {
		t.Errorf("expected status id to be %d, was %d", expected.id, status.id)
	}
	if !status.isActive {
		t.Errorf("expected archetype status to be 'true', was %t", status.isActive)
	}
}

func TestScene_AddComponent_WhenAddingMultipleComponents_ShouldUpdateArchetype(t *testing.T) {
	comp1 := TransformComponent{mgl.Vec3{1, 2, 3}, mgl.QuatIdent()}
	comp2 := MeshComponent{}
	scene := NewScene(200)
	entity := scene.NewEntity()

	// Act
	scene.AddTfComp(entity, comp1)
	arch1 := scene.archetypes.archEntities[entity.index]
	scene.AddMeshComp(entity, comp2)
	arch2 := scene.archetypes.archEntities[entity.index]

	// Assert
	if arch1.id == arch2.id {
		t.Errorf("expected archetype IDs to be different, were same %d", arch1.id)
	}
}

func TestContains(t *testing.T) {
	entities := []uint32{42, 2, 8}
	if !util.Contains(entities, 42) {
		t.Errorf("expected 42 to be in %v", entities)
	}
	if util.Contains(entities, 82) {
		t.Errorf("Did not expect to find 82 in %v", entities)
	}
}

func TestEntityIDAllocator_SingleAllocation_ShouldReturn0(t *testing.T) {
	alloc := newEntityAllocator(10)
	actual := alloc.allocate()
	expected := Entity{index: 0, version: 1}
	if expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestEntityIDAllocator_TwoAllocations_ShouldIncrementIndex(t *testing.T) {
	alloc := newEntityAllocator(10)
	alloc.allocate()
	actual := alloc.allocate()
	expected := Entity{index: 1, version: 1}
	if expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestEntityIDAllocator_Free_ShouldIncreaseVersionAndMoveToFree(t *testing.T) {
	alloc := newEntityAllocator(10)
	entity := alloc.allocate()
	if !alloc.entities[entity.index].isActive {
		t.Errorf("expected entity to be active, was %v", alloc.entities[entity.index])
	}
	alloc.deallocate(entity)

	if alloc.entities[0].version != 1 {
		t.Error("expected entity version 0")
	}
	if len(alloc.free) != 10 {
		t.Errorf("expected 10 free indices, have %d", len(alloc.free))
	}
	if alloc.entities[entity.index].isActive {
		t.Errorf("expected entity %v to be inactive, was %v", entity, alloc.entities[entity.index])
	}
}
