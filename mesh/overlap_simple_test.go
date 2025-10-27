package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

// TestSimpleOverlapCheck is a minimal test to verify WithTriangleOverlapCheck works
func TestSimpleOverlapCheck(t *testing.T) {
	m := NewMesh(
		WithEpsilon(1e-9),
		WithTriangleOverlapCheck(true), // Enable overlap checking
	)

	// Add first triangle
	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 5, Y: 10})

	err := m.AddTriangle(v0, v1, v2)
	if err != nil {
		t.Fatalf("Failed to add first triangle: %v", err)
	}

	t.Logf("First triangle added successfully")
	t.Logf("Mesh now has %d triangles", m.NumTriangles())

	// Try to add second triangle that overlaps
	v3, _ := m.AddVertex(types.Point{X: 3, Y: 3})
	v4, _ := m.AddVertex(types.Point{X: 7, Y: 3})
	v5, _ := m.AddVertex(types.Point{X: 5, Y: 7})

	err = m.AddTriangle(v3, v4, v5)

	t.Logf("Attempting to add overlapping triangle...")
	t.Logf("Second triangle add result: err=%v", err)

	if err == nil {
		t.Errorf("Expected error when adding overlapping triangle, but got none")
		t.Logf("Mesh now has %d triangles", m.NumTriangles())

		// Check if FindOverlappingTriangles detects it
		overlaps := m.FindOverlappingTriangles()
		t.Logf("FindOverlappingTriangles() found %d overlaps", len(overlaps))
		if len(overlaps) > 0 {
			for i, overlap := range overlaps {
				t.Logf("Overlap %d: area=%.4f, type=%s", i+1, overlap.IntersectionArea, overlap.Type)
			}
		}
	} else {
		t.Logf("✓ Second triangle correctly rejected: %v", err)
	}
}
