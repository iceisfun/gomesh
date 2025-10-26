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

func TestGetUntriangulatedVertices(t *testing.T) {
	m := NewMesh()

	// Create a square with 4 vertices
	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 0, Y: 10})

	// Create a second square with 4 more vertices (sharing no vertices with first)
	v4, _ := m.AddVertex(types.Point{X: 20, Y: 0})
	v5, _ := m.AddVertex(types.Point{X: 30, Y: 0})
	v6, _ := m.AddVertex(types.Point{X: 30, Y: 10})
	v7, _ := m.AddVertex(types.Point{X: 20, Y: 10})

	// Only triangulate the first square
	_ = m.AddTriangle(v0, v1, v2)
	_ = m.AddTriangle(v0, v2, v3)

	// Create loops for both squares
	loop1 := types.NewPolygonLoop(v0, v1, v2, v3)
	loop2 := types.NewPolygonLoop(v4, v5, v6, v7)

	t.Run("all vertices triangulated", func(t *testing.T) {
		untriangulated := m.GetUntriangulatedVertices([]types.PolygonLoop{loop1})
		if len(untriangulated) != 0 {
			t.Errorf("Expected 0 untriangulated vertices in loop1, got %d", len(untriangulated))
		}
	})

	t.Run("no vertices triangulated", func(t *testing.T) {
		untriangulated := m.GetUntriangulatedVertices([]types.PolygonLoop{loop2})
		if len(untriangulated) != 4 {
			t.Errorf("Expected 4 untriangulated vertices in loop2, got %d", len(untriangulated))
		}
		// Verify all vertices from loop2 are in the result
		untriangulatedSet := make(map[types.VertexID]bool)
		for _, vid := range untriangulated {
			untriangulatedSet[vid] = true
		}
		for _, vid := range []types.VertexID{v4, v5, v6, v7} {
			if !untriangulatedSet[vid] {
				t.Errorf("Expected vertex %d to be in untriangulated set", vid)
			}
		}
	})

	t.Run("mixed triangulation", func(t *testing.T) {
		untriangulated := m.GetUntriangulatedVertices([]types.PolygonLoop{loop1, loop2})
		if len(untriangulated) != 4 {
			t.Errorf("Expected 4 untriangulated vertices total, got %d", len(untriangulated))
		}
	})

	t.Run("empty loops", func(t *testing.T) {
		untriangulated := m.GetUntriangulatedVertices([]types.PolygonLoop{})
		if len(untriangulated) != 0 {
			t.Errorf("Expected 0 untriangulated vertices for empty loops, got %d", len(untriangulated))
		}
	})

	t.Run("partial triangulation", func(t *testing.T) {
		m2 := NewMesh()
		p0, _ := m2.AddVertex(types.Point{X: 0, Y: 0})
		p1, _ := m2.AddVertex(types.Point{X: 10, Y: 0})
		p2, _ := m2.AddVertex(types.Point{X: 10, Y: 10})
		p3, _ := m2.AddVertex(types.Point{X: 0, Y: 10})

		// Only triangulate 3 of the 4 vertices
		_ = m2.AddTriangle(p0, p1, p2)

		loop := types.NewPolygonLoop(p0, p1, p2, p3)
		untriangulated := m2.GetUntriangulatedVertices([]types.PolygonLoop{loop})

		if len(untriangulated) != 1 {
			t.Errorf("Expected 1 untriangulated vertex, got %d", len(untriangulated))
		}
		if len(untriangulated) > 0 && untriangulated[0] != p3 {
			t.Errorf("Expected vertex %d to be untriangulated, got %d", p3, untriangulated[0])
		}
	})
}
