package ecs

import (
	"testing"
)

func TestMakeComponentMask(t *testing.T) {
    componentA := ComponentID(3)  // binary: 0011
    componentB := ComponentID(10) // binary: 1010
    expected := ComponentMask(11) // binary: 1011
    actual := makeComponentMask(componentA, componentB)

    if expected != actual {
        t.Errorf("expected: %d, actual: %d", expected, actual)
    }
}

func TestNewComponentMatrix_ShouldAddEntitiesAndComponents(t *testing.T) {
    // Arrange
    nEntities, nComponents := uint(3), uint(16)

    // Act
    mat := newComponentMatrix(nEntities, nComponents)

    // Assert
    rows := len(mat)
    cols := mat[0]
    if uint(rows) != nEntities {
        t.Errorf("expected number of rows to be %d", nEntities)
    }
    if cols != 0 {
        t.Errorf("expected component mask to be 0, was %d", cols)
    }
}

func TestEntityComponents_newEntityComponents_ShouldAddEntities(t *testing.T) {
    maxEntities := uint(3)
    entities := newEntities(maxEntities)
    if entities.nEntities != maxEntities {
        t.Errorf("expected %d number of entities, was %d", maxEntities, entities.nEntities)
    }
    lenAllocEntities := len(entities.allocator.entities)
    lenAllocFree := len(entities.allocator.free)
    if lenAllocEntities != lenAllocFree || uint(lenAllocEntities) != maxEntities {
        t.Errorf("number of allocator entities was not equal to number of free indices, was no equal to %d", maxEntities)
    }
    lenCompMat := len(entities.matrix)
    if uint(lenCompMat) != maxEntities {
        t.Errorf("expected len comp matrix to be %d", maxEntities)
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
    lenCompMat := len(scene.entities.matrix)
    if uint(lenCompMat) != maxEntities {
        t.Errorf("expected len comp matrix to be %d", maxEntities)
    }
}

func TestScene_AddEntity_ShouldBe0(t *testing.T) {
    maxEntities := uint(3)
    scene := NewScene(maxEntities)
    entity := scene.NewEntity()
    actual := EntityID{0, 1}
    if entity != actual {
        t.Errorf("expected allocated entity to be %v, was %v", entity, actual)
    }
}

func TestContains(t *testing.T) {
    entities := []uint32{42, 2, 8}
    if !contains(entities, 42) {
        t.Errorf("expected 42 to be in %v", entities)
    }
    if contains(entities, 82) {
        t.Errorf("Did not expect to find 82 in %v", entities)
    }
}

func TestEntityIDAllocator_SingleAllocation_ShouldReturn0(t *testing.T) {
    alloc := newEntityAllocator(10)
    actual := alloc.allocate()
    expected := EntityID{index: 0, version: 1}
    if expected != actual {
        t.Errorf("expected %v, actual %v", expected, actual)
    }
}

func TestEntityIDAllocator_TwoAllocations_ShouldIncrementIndex(t *testing.T) {
    alloc := newEntityAllocator(10)
    alloc.allocate()
    actual := alloc.allocate()
    expected := EntityID{index: 1, version: 1}
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
