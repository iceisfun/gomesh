package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestMeshGetters(t *testing.T) {
	m := NewMesh()
	a, _ := m.AddVertex(types.Point{1, 0})
	b, _ := m.AddVertex(types.Point{2, 0})
	c, _ := m.AddVertex(types.Point{0, 2})
	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.NumVertices() != 3 {
		t.Fatalf("expected 3 vertices")
	}
	if m.NumTriangles() != 1 {
		t.Fatalf("expected 1 triangle")
	}

	verts := m.GetVertices()
	if len(verts) != 3 {
		t.Fatalf("expected copy of vertices")
	}
	verts[0] = types.Point{}
	orig := m.GetVertex(a)
	if orig.X == 0 && orig.Y == 0 {
		t.Fatalf("expected defensive copy of vertices")
	}

	ts := m.GetTriangles()
	if len(ts) != 1 {
		t.Fatalf("expected copy of triangles")
	}
}
