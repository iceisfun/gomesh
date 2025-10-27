package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

// TestEdgeReuseValidation tests that we cannot add a third triangle to an edge
// that already has 2 triangles.
func TestEdgeReuseValidation(t *testing.T) {
	// Use a simpler configuration that only enables edge intersection check
	m := NewMesh(
		WithEdgeIntersectionCheck(true),
	)

	// Create vertices for 3 triangles sharing an edge
	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{10, 0})
	v2, _ := m.AddVertex(types.Point{5, 5})
	v3, _ := m.AddVertex(types.Point{5, -5})
	v4, _ := m.AddVertex(types.Point{5, 2}) // Not collinear with v0-v1

	// Add first triangle using edge v0-v1
	if err := m.AddTriangle(v0, v1, v2); err != nil {
		t.Fatalf("failed to add first triangle: %v", err)
	}

	// Add second triangle using edge v0-v1 (should succeed)
	if err := m.AddTriangle(v1, v0, v3); err != nil {
		t.Fatalf("failed to add second triangle: %v", err)
	}

	// Try to add third triangle using edge v0-v1 (should FAIL - edge already has 2 triangles)
	err := m.AddTriangle(v0, v1, v4)
	if err == nil {
		t.Error("Expected error when adding third triangle to edge v0-v1, got nil")
	} else {
		t.Logf("Correctly rejected third triangle: %v", err)
	}
}

// TestEdgeReuseProblemFromArea1 reproduces the specific edge reuse from area_1.json
func TestEdgeReuseProblemFromArea1(t *testing.T) {
	m := NewMesh(
		WithEdgeIntersectionCheck(true),
		WithOverlapTriangle(false),
	)

	// Recreate the edge [16-25] scenario from area_1.json
	// Edge [16-25] at (151,146) to (150,154) is shared by 4 triangles

	v15, _ := m.AddVertex(types.Point{151, 144})
	v16, _ := m.AddVertex(types.Point{151, 146})
	v19, _ := m.AddVertex(types.Point{149, 150})
	v20, _ := m.AddVertex(types.Point{151, 152})
	v24, _ := m.AddVertex(types.Point{149, 154})
	v25, _ := m.AddVertex(types.Point{150, 154})

	// Try to add 4 triangles all sharing edge 16-25
	// First triangle should succeed
	if err := m.AddTriangle(v20, v16, v25); err != nil {
		t.Fatalf("failed to add first triangle: %v", err)
	}

	// Second triangle should succeed
	if err := m.AddTriangle(v15, v16, v25); err != nil {
		t.Fatalf("failed to add second triangle: %v", err)
	}

	// Third triangle should FAIL (edge 16-25 already has 2 triangles)
	err := m.AddTriangle(v25, v16, v24)
	if err == nil {
		t.Error("Expected error when adding third triangle to edge 16-25, got nil")
	} else {
		t.Logf("Correctly rejected third triangle: %v", err)
	}

	// Fourth triangle should also FAIL
	err = m.AddTriangle(v16, v19, v25)
	if err == nil {
		t.Error("Expected error when adding fourth triangle to edge 16-25, got nil")
	} else {
		t.Logf("Correctly rejected fourth triangle: %v", err)
	}
}
