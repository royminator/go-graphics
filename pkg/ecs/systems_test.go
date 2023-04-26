package ecs

import (
    "testing"
)

func TestMakeComponentMask(t *testing.T) {
    componentA := ComponentID(2)  // binary: 0010
    componentB := ComponentID(8)  // binary: 1000
    expected := ComponentMask(10) // binary: 1010
    actual := makeComponentMask(componentA, componentB)

    if expected != actual {
        t.Errorf("expected: %d, actual: %d", expected, actual)
    }
}
