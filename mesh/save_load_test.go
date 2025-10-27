package mesh

import (
	"os"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestSaveLoad(t *testing.T) {
	// Create a mesh with some configuration
	m := NewMesh(
		WithEpsilon(1e-9),
		WithMergeVertices(true),
		WithEdgeIntersectionCheck(true),
		WithTriangleEnforceNoVertexInside(true),
		WithEdgeCannotCrossPerimeter(true),
		WithOverlapTriangle(false),
	)

	// Add a perimeter
	perimeter := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}
	_, err := m.AddPerimeter(perimeter)
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add a hole
	hole := []types.Point{
		{X: 3, Y: 3},
		{X: 7, Y: 3},
		{X: 7, Y: 7},
		{X: 3, Y: 7},
	}
	_, err = m.AddHole(hole)
	if err != nil {
		t.Fatalf("failed to add hole: %v", err)
	}

	// Add some triangles
	v1, _ := m.AddVertex(types.Point{2, 2})
	v2, _ := m.AddVertex(types.Point{8, 2})
	v3, _ := m.AddVertex(types.Point{5, 8})

	if err := m.AddTriangle(0, 1, v1); err != nil {
		t.Fatalf("failed to add triangle: %v", err)
	}

	// Save to temp file
	tmpfile := "/tmp/test_mesh.json"
	defer os.Remove(tmpfile)

	if err := m.Save(tmpfile); err != nil {
		t.Fatalf("failed to save mesh: %v", err)
	}

	// Load the mesh back
	m2, err := Load(tmpfile)
	if err != nil {
		t.Fatalf("failed to load mesh: %v", err)
	}

	// Verify the loaded mesh matches
	if m2.NumVertices() != m.NumVertices() {
		t.Errorf("vertex count mismatch: got %d, want %d", m2.NumVertices(), m.NumVertices())
	}

	if m2.NumTriangles() != m.NumTriangles() {
		t.Errorf("triangle count mismatch: got %d, want %d", m2.NumTriangles(), m.NumTriangles())
	}

	if len(m2.perimeters) != len(m.perimeters) {
		t.Errorf("perimeter count mismatch: got %d, want %d", len(m2.perimeters), len(m.perimeters))
	}

	if len(m2.holes) != len(m.holes) {
		t.Errorf("hole count mismatch: got %d, want %d", len(m2.holes), len(m.holes))
	}

	// Verify config was preserved
	if m2.cfg.epsilon != m.cfg.epsilon {
		t.Errorf("epsilon mismatch: got %v, want %v", m2.cfg.epsilon, m.cfg.epsilon)
	}

	if m2.cfg.validateEdgeCannotCrossPerimeter != m.cfg.validateEdgeCannotCrossPerimeter {
		t.Errorf("validateEdgeCannotCrossPerimeter mismatch")
	}

	if m2.cfg.errorOnDuplicateTriangle != m.cfg.errorOnDuplicateTriangle {
		t.Errorf("errorOnDuplicateTriangle mismatch")
	}

	t.Logf("Successfully saved and loaded mesh with %d vertices, %d triangles", m2.NumVertices(), m2.NumTriangles())

	// Verify we can still use the loaded mesh
	_, _ = v2, v3
}

func TestSaveLoadPreservesGeometry(t *testing.T) {
	// Create a simple mesh
	m := NewMesh()
	v0, _ := m.AddVertex(types.Point{1.5, 2.5})
	v1, _ := m.AddVertex(types.Point{5.5, 3.5})
	v2, _ := m.AddVertex(types.Point{3.5, 7.5})

	if err := m.AddTriangle(v0, v1, v2); err != nil {
		t.Fatalf("failed to add triangle: %v", err)
	}

	// Save and load
	tmpfile := "/tmp/test_mesh_geometry.json"
	defer os.Remove(tmpfile)

	if err := m.Save(tmpfile); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	m2, err := Load(tmpfile)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	// Verify vertices match exactly
	for i := 0; i < m.NumVertices(); i++ {
		p1 := m.vertices[i]
		p2 := m2.vertices[i]
		if p1.X != p2.X || p1.Y != p2.Y {
			t.Errorf("vertex %d mismatch: got %v, want %v", i, p2, p1)
		}
	}

	// Verify triangles match
	if len(m.triangles) != len(m2.triangles) {
		t.Fatalf("triangle count mismatch")
	}

	for i, tri := range m.triangles {
		tri2 := m2.triangles[i]
		if tri.V1() != tri2.V1() || tri.V2() != tri2.V2() || tri.V3() != tri2.V3() {
			t.Errorf("triangle %d mismatch: got %v, want %v", i, tri2, tri)
		}
	}
}
