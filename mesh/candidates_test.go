package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestVertexFindCandidates(t *testing.T) {
	m := NewMesh()

	// Create a simple scenario
	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{10, 0})
	v2, _ := m.AddVertex(types.Point{10, 10})
	v3, _ := m.AddVertex(types.Point{0, 10})

	// Add a perimeter
	_, err := m.AddPerimeter([]types.Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add interior vertex
	v4, _ := m.AddVertex(types.Point{5, 5})

	// Add exterior vertex (should not be connectable)
	v5, _ := m.AddVertex(types.Point{15, 5})

	// Find candidates for interior vertex
	candidates := m.VertexFindCandidates(v4)

	// v4 should be able to connect to v0, v1, v2, v3 (perimeter vertices)
	// but not v5 (would cross perimeter)
	if len(candidates) < 4 {
		t.Errorf("Expected at least 4 candidates, got %d", len(candidates))
	}

	// Check that v5 is not in candidates
	for _, c := range candidates {
		if c.VertexID == v5 {
			t.Error("Exterior vertex should not be a candidate")
		}
	}

	_ = v0
	_ = v1
	_ = v2
	_ = v3
}

func TestVertexFindCandidatesWithConstraints(t *testing.T) {
	m := NewMesh(WithEdgeCannotCrossPerimeter(true))

	// Create perimeter
	_, _ = m.AddVertex(types.Point{0, 0})
	_, _ = m.AddVertex(types.Point{20, 0})
	_, _ = m.AddVertex(types.Point{20, 20})
	_, _ = m.AddVertex(types.Point{0, 20})

	_, err := m.AddPerimeter([]types.Point{{0, 0}, {20, 0}, {20, 20}, {0, 20}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add interior vertices
	v1, _ := m.AddVertex(types.Point{5, 10})
	v2, _ := m.AddVertex(types.Point{15, 10})

	// Add triangle creating an edge
	v3, _ := m.AddVertex(types.Point{10, 5})
	_ = m.AddTriangle(v1, v2, v3)

	// Find candidates for a new vertex
	v4, _ := m.AddVertex(types.Point{10, 15})
	candidates := m.VertexFindCandidates(v4)

	// Should find some candidates
	if len(candidates) == 0 {
		t.Error("Expected some candidates")
	}

	t.Logf("Found %d candidates for vertex %d", len(candidates), v4)
}

func TestVertexFindTriangleCandidates(t *testing.T) {
	m := NewMesh()

	// Create a simple scenario with 4 vertices
	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{10, 0})
	v2, _ := m.AddVertex(types.Point{10, 10})
	v3, _ := m.AddVertex(types.Point{0, 10})

	// Find triangle candidates for v0
	candidates := m.VertexFindTriangleCandidates(v0)

	// With 4 vertices, v0 can form triangles with combinations of the other 3
	// Should find at least a few valid triangles
	if len(candidates) == 0 {
		t.Error("Expected some triangle candidates")
	}

	t.Logf("Found %d triangle candidates for vertex %d", len(candidates), v0)

	// Verify each candidate includes v0
	for _, tri := range candidates {
		if tri.V1 != v0 && tri.V2 != v0 && tri.V3 != v0 {
			t.Errorf("Triangle candidate doesn't include vertex %d", v0)
		}
	}

	_ = v1
	_ = v2
	_ = v3
}

func TestVertexFindTriangleCandidatesWithPerimeter(t *testing.T) {
	m := NewMesh(WithEdgeCannotCrossPerimeter(true))

	// Create perimeter
	_, _ = m.AddVertex(types.Point{0, 0})
	_, _ = m.AddVertex(types.Point{10, 0})
	_, _ = m.AddVertex(types.Point{10, 10})
	_, _ = m.AddVertex(types.Point{0, 10})

	_, err := m.AddPerimeter([]types.Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add interior vertex
	v1, _ := m.AddVertex(types.Point{2, 2})
	v2, _ := m.AddVertex(types.Point{8, 2})
	v3, _ := m.AddVertex(types.Point{5, 8})

	// Add exterior vertex
	_, _ = m.AddVertex(types.Point{15, 5})

	// Find triangle candidates for v1
	candidates := m.VertexFindTriangleCandidates(v1)

	// Should find some interior triangles
	if len(candidates) == 0 {
		t.Error("Expected some triangle candidates")
	}

	t.Logf("Found %d triangle candidates with perimeter constraint", len(candidates))

	// None of the candidates should cross the perimeter
	// (this is enforced by the validation)
	for _, tri := range candidates {
		// Try to add it (should succeed)
		testMesh := NewMesh(WithEdgeCannotCrossPerimeter(true))
		testMesh.vertices = m.vertices
		testMesh.perimeters = m.perimeters

		err := testMesh.AddTriangle(tri.V1, tri.V2, tri.V3)
		if err != nil {
			t.Errorf("Candidate triangle should be valid: %v", err)
		}
	}

	_ = v2
	_ = v3
}

func TestVertexFindCandidatesInvalidVertex(t *testing.T) {
	m := NewMesh()

	// Test with invalid vertex ID
	candidates := m.VertexFindCandidates(999)
	if candidates != nil {
		t.Error("Expected nil for invalid vertex")
	}
}

func TestVertexFindTriangleCandidatesInvalidVertex(t *testing.T) {
	m := NewMesh()

	// Test with invalid vertex ID
	candidates := m.VertexFindTriangleCandidates(999)
	if candidates != nil {
		t.Error("Expected nil for invalid vertex")
	}
}

func TestVertexFindCandidatesConcavePerimeter(t *testing.T) {
	m := NewMesh()

	// Create a concave "C" shape perimeter (like pac-man)
	// The opening faces right
	/*
	   0,10 ---- 10,10
	     |          |
	     |    7,7   |
	     |          |
	   0,3         10,3
	     |          |
	   0,0 ---- 10,0
	*/

	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{10, 0})
	v2, _ := m.AddVertex(types.Point{10, 3})
	v3, _ := m.AddVertex(types.Point{7, 3})
	v4, _ := m.AddVertex(types.Point{7, 7})
	v5, _ := m.AddVertex(types.Point{10, 7})
	v6, _ := m.AddVertex(types.Point{10, 10})
	v7, _ := m.AddVertex(types.Point{0, 10})

	_, err := m.AddPerimeter([]types.Point{
		{0, 0}, {10, 0}, {10, 3}, {7, 3}, {7, 7}, {10, 7}, {10, 10}, {0, 10},
	})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add interior vertex
	v8, _ := m.AddVertex(types.Point{5, 5})

	// Find candidates from bottom-left corner (v0)
	candidates := m.VertexFindCandidates(v0)

	// v0 should be able to connect to interior and perimeter vertices
	// but NOT to vertices on the opposite side of the concave opening
	// if the edge would go outside the perimeter

	// Check that v3 (left side of opening) is reachable from v0
	foundV3 := false
	for _, c := range candidates {
		if c.VertexID == v3 {
			foundV3 = true
		}
	}
	if !foundV3 {
		t.Error("Expected v3 to be a candidate from v0")
	}

	// The edge from v7 (top-left) to v2 (bottom-right of opening)
	// should be rejected because it would go outside the concave perimeter
	candidatesFromV7 := m.VertexFindCandidates(v7)
	foundV2FromV7 := false
	for _, c := range candidatesFromV7 {
		if c.VertexID == v2 {
			foundV2FromV7 = true
		}
	}
	// This edge's midpoint would be at (5, 6.5) which should be outside
	// the perimeter due to the concave opening
	if foundV2FromV7 {
		t.Error("Expected v2 to NOT be a candidate from v7 (edge goes outside concave perimeter)")
	}

	t.Logf("Found %d candidates from v0", len(candidates))
	t.Logf("Found %d candidates from v7", len(candidatesFromV7))

	_ = v1
	_ = v4
	_ = v5
	_ = v6
	_ = v8
}
