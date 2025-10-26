package mesh

import (
	"testing"

	"gomesh/types"
)

func TestAddVertexWithoutMerging(t *testing.T) {
	m := NewMesh()
	id1, err := m.AddVertex(types.Point{X: 0, Y: 0})
	if err != nil || id1 != 0 {
		t.Fatalf("unexpected result: id=%v err=%v", id1, err)
	}
	id2, err := m.AddVertex(types.Point{X: 0, Y: 0})
	if err != nil || id2 == id1 {
		t.Fatalf("expected new vertex when merging disabled")
	}
}

func TestAddVertexWithMerging(t *testing.T) {
	m := NewMesh(WithMergeVertices(true), WithMergeDistance(0.5))
	id1, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	id2, _ := m.AddVertex(types.Point{X: 0.1, Y: 0.1})
	if id1 != id2 {
		t.Fatalf("expected merge to reuse vertex id")
	}

	if m.NumVertices() != 1 {
		t.Fatalf("expected single stored vertex")
	}
}

func TestFindVertexNear(t *testing.T) {
	m := NewMesh(WithMergeVertices(true), WithMergeDistance(0.5))
	id, _ := m.AddVertex(types.Point{X: 1, Y: 1})
	found, ok := m.FindVertexNear(types.Point{X: 1.1, Y: 1.1})
	if !ok || found != id {
		t.Fatalf("expected to locate nearby vertex")
	}
}
