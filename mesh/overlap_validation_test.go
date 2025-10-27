package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

// TestOverlap_T414_T419 tests a real overlap case found in area_1.json
// This test documents a VALIDATION BUG where overlapping triangles are not rejected.
func TestOverlap_T414_T419(t *testing.T) {
	m := NewMesh(
		WithEpsilon(1e-9),
		WithMergeVertices(true),
		WithEdgeIntersectionCheck(true),
		WithTriangleEnforceNoVertexInside(true),
		WithEdgeCannotCrossPerimeter(true),
		WithOverlapTriangle(false),
		WithTriangleOverlapCheck(true), // Enable geometric overlap validation
	)

	// Add vertices
	v0, err := m.AddVertex(types.Point{X: 71.00, Y: 150.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 0: %v", err)
	}
	v1, err := m.AddVertex(types.Point{X: 62.00, Y: 127.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 1: %v", err)
	}
	v2, err := m.AddVertex(types.Point{X: 72.00, Y: 150.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 2: %v", err)
	}
	v3, err := m.AddVertex(types.Point{X: 62.00, Y: 150.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 3: %v", err)
	}
	v4, err := m.AddVertex(types.Point{X: 62.00, Y: 125.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 4: %v", err)
	}

	// Add first triangle
	err = m.AddTriangle(v0, v1, v2)
	if err != nil {
		t.Fatalf("Failed to add first triangle: %v", err)
	}

	// Add second triangle (should fail due to overlap)
	// These triangles have an intersection area of 11.5
	err = m.AddTriangle(v3, v4, v2)
	if err == nil {
		t.Errorf("Expected error when adding overlapping triangle (intersection area: 11.5000), but got none")

		// Print debug information
		t.Logf("Triangle 1: v0=%v, v1=%v, v2=%v", m.vertices[v0], m.vertices[v1], m.vertices[v2])
		t.Logf("Triangle 2: v3=%v, v4=%v, v2=%v", m.vertices[v3], m.vertices[v4], m.vertices[v2])

		// Check for actual overlap
		overlaps := m.FindOverlappingTriangles()
		t.Logf("Found %d overlaps", len(overlaps))
		for i, overlap := range overlaps {
			t.Logf("Overlap %d: area=%.4f, type=%s", i+1, overlap.IntersectionArea, overlap.Type)
		}
	} else {
		t.Logf("✓ Second triangle correctly rejected: %v", err)
	}
}

// TestOverlap_T16_T17 tests another real overlap case from area_1_example_2.json
func TestOverlap_T16_T17(t *testing.T) {
	m := NewMesh(
		WithEpsilon(1e-9),
		WithMergeVertices(true),
		WithEdgeIntersectionCheck(true),
		WithTriangleEnforceNoVertexInside(true),
		WithEdgeCannotCrossPerimeter(true),
		WithOverlapTriangle(false),
		WithTriangleOverlapCheck(true), // Enable geometric overlap validation
	)

	// Add vertices
	v0, err := m.AddVertex(types.Point{X: 90.00, Y: 125.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 0: %v", err)
	}
	v1, err := m.AddVertex(types.Point{X: 87.00, Y: 165.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 1: %v", err)
	}
	v2, err := m.AddVertex(types.Point{X: 87.00, Y: 135.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 2: %v", err)
	}
	v3, err := m.AddVertex(types.Point{X: 87.00, Y: 166.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 3: %v", err)
	}
	v4, err := m.AddVertex(types.Point{X: 87.00, Y: 136.00})
	if err != nil {
		t.Fatalf("Failed to add vertex 4: %v", err)
	}

	// Add first triangle
	err = m.AddTriangle(v0, v1, v2)
	if err != nil {
		t.Fatalf("Failed to add first triangle: %v", err)
	}

	// Add second triangle (should fail due to overlap)
	// These triangles have an intersection area of 43.5
	err = m.AddTriangle(v0, v3, v4)
	if err == nil {
		t.Errorf("Expected error when adding overlapping triangle (intersection area: 43.5000), but got none")

		// Print debug information
		t.Logf("Triangle 1: v0=%v, v1=%v, v2=%v", m.vertices[v0], m.vertices[v1], m.vertices[v2])
		t.Logf("Triangle 2: v0=%v, v3=%v, v4=%v", m.vertices[v0], m.vertices[v3], m.vertices[v4])

		// Check for actual overlap
		overlaps := m.FindOverlappingTriangles()
		t.Logf("Found %d overlaps", len(overlaps))
		for i, overlap := range overlaps {
			t.Logf("Overlap %d: area=%.4f, type=%s", i+1, overlap.IntersectionArea, overlap.Type)
		}
	} else {
		t.Logf("✓ Second triangle correctly rejected: %v", err)
	}
}
