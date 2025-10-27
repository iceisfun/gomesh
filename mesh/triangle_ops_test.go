package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestAddTriangleSuccess(t *testing.T) {
	m := NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{1, 0})
	c, _ := m.AddVertex(types.Point{0, 1})

	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("unexpected error adding triangle: %v", err)
	}
	if m.NumTriangles() != 1 {
		t.Fatalf("expected triangle count 1")
	}
}

func TestAddTriangleInvalidVertex(t *testing.T) {
	m := NewMesh()
	if err := m.AddTriangle(0, 1, 2); err != ErrInvalidVertexID {
		t.Fatalf("expected invalid vertex error, got %v", err)
	}
}

func TestAddTriangleDegenerate(t *testing.T) {
	m := NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{1, 1})
	c, _ := m.AddVertex(types.Point{2, 2})

	if err := m.AddTriangle(a, b, c); err != ErrDegenerateTriangle {
		t.Fatalf("expected degenerate triangle error, got %v", err)
	}
}

func TestDuplicateTriangleError(t *testing.T) {
	m := NewMesh(WithDuplicateTriangleError(true))
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{1, 0})
	c, _ := m.AddVertex(types.Point{0, 1})

	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.AddTriangle(b, c, a); err != ErrDuplicateTriangle {
		t.Fatalf("expected duplicate triangle error, got %v", err)
	}
}

func TestOpposingWindingDuplicate(t *testing.T) {
	m := NewMesh(WithDuplicateTriangleOpposingWinding(true))
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{1, 0})
	c, _ := m.AddVertex(types.Point{0, 1})

	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("same winding should be allowed, got %v", err)
	}
	if err := m.AddTriangle(a, c, b); err != ErrOpposingWindingDuplicate {
		t.Fatalf("expected opposing winding error, got %v", err)
	}
}

func TestVertexInsideValidation(t *testing.T) {
	m := NewMesh(WithTriangleEnforceNoVertexInside(true))
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{4, 0})
	c, _ := m.AddVertex(types.Point{0, 4})

	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("unexpected error adding first triangle: %v", err)
	}

	if _, err := m.AddVertex(types.Point{1, 0.5}); err != nil {
		t.Fatalf("failed to add inner vertex: %v", err)
	}
	d, err := m.AddVertex(types.Point{4, 4})
	if err != nil {
		t.Fatalf("failed to add outer vertex: %v", err)
	}

	if err := m.AddTriangle(a, b, d); err != ErrVertexInsideTriangle {
		t.Fatalf("expected vertex-inside error, got %v", err)
	}
}

func TestEdgeIntersectionValidation(t *testing.T) {
	m := NewMesh(WithEdgeIntersectionCheck(true))

	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{2, 0})
	v2, _ := m.AddVertex(types.Point{0, 2})
	v3, _ := m.AddVertex(types.Point{1, -1})
	v4, _ := m.AddVertex(types.Point{1, 1})

	if err := m.AddTriangle(v0, v1, v2); err != nil {
		t.Fatalf("unexpected error adding first triangle: %v", err)
	}

	if err := m.AddTriangle(v3, v1, v4); err != ErrEdgeIntersection {
		t.Fatalf("expected edge intersection error, got %v", err)
	}
}

func TestEdgeCannotCrossPerimeter(t *testing.T) {
	m := NewMesh(WithEdgeCannotCrossPerimeter(true))

	// Create a square perimeter
	p0, _ := m.AddVertex(types.Point{0, 0})
	p1, _ := m.AddVertex(types.Point{10, 0})
	_, _ = m.AddVertex(types.Point{10, 10})
	_, _ = m.AddVertex(types.Point{0, 10})

	_, err := m.AddPerimeter([]types.Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add interior vertices
	v1, _ := m.AddVertex(types.Point{2, 2})
	v2, _ := m.AddVertex(types.Point{8, 2})
	v3, _ := m.AddVertex(types.Point{5, 8})

	// Triangle with edge along perimeter edge (should succeed)
	if err := m.AddTriangle(p0, p1, v1); err != nil {
		t.Fatalf("triangle with edge on perimeter should succeed, got %v", err)
	}

	// Triangle with edge crossing perimeter (should fail)
	// Edge from v1 to v2 stays inside, but trying to connect to outside
	v4, _ := m.AddVertex(types.Point{15, 5})
	if err := m.AddTriangle(v1, v2, v4); err != ErrEdgeCrossesPerimeter {
		t.Fatalf("expected perimeter crossing error, got %v", err)
	}

	// Triangle entirely inside perimeter (should succeed)
	if err := m.AddTriangle(v1, v2, v3); err != nil {
		t.Fatalf("triangle inside perimeter should succeed, got %v", err)
	}
}

func TestEdgeCannotCrossHole(t *testing.T) {
	m := NewMesh(
		WithEdgeCannotCrossPerimeter(true),
		WithMergeVertices(true),
	)

	// Create outer perimeter
	_, err := m.AddPerimeter([]types.Point{{0, 0}, {20, 0}, {20, 20}, {0, 20}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Create hole in the middle
	h0, _ := m.AddVertex(types.Point{8, 8})
	h1, _ := m.AddVertex(types.Point{12, 8})
	_, _ = m.AddVertex(types.Point{12, 12})
	_, _ = m.AddVertex(types.Point{8, 12})

	_, err = m.AddHole([]types.Point{{8, 8}, {12, 8}, {12, 12}, {8, 12}})
	if err != nil {
		t.Fatalf("failed to add hole: %v", err)
	}

	// Add vertices below the hole
	v1, _ := m.AddVertex(types.Point{10, 5}) // Below hole, centered
	v2, _ := m.AddVertex(types.Point{15, 10}) // Right of hole

	// Triangle with edge along hole boundary (should succeed)
	// Uses the bottom edge of the hole (h0->h1) which is at y=8
	if err := m.AddTriangle(h0, h1, v1); err != nil {
		t.Fatalf("triangle with edge on hole boundary should succeed, got %v", err)
	}

	// Triangle with edge crossing hole boundary (should fail)
	// Try to connect vertex inside hole region to vertex outside
	v3, _ := m.AddVertex(types.Point{10, 10}) // Inside hole
	if err := m.AddTriangle(v1, v2, v3); err != ErrEdgeCrossesPerimeter {
		t.Fatalf("expected hole crossing error, got %v", err)
	}
}

func TestOverlapTriangleDefault(t *testing.T) {
	// By default, overlapping triangles are allowed
	m := NewMesh()

	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{10, 0})
	v2, _ := m.AddVertex(types.Point{5, 10})

	_, err := m.AddPerimeter([]types.Point{{0, 0}, {10, 0}, {5, 10}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add triangle with one vertex order
	if err := m.AddTriangle(v0, v1, v2); err != nil {
		t.Fatalf("unexpected error adding first triangle: %v", err)
	}

	// Add same triangle with rotated vertex order (should succeed by default)
	if err := m.AddTriangle(v1, v2, v0); err != nil {
		t.Fatalf("expected overlapping triangle to succeed by default, got %v", err)
	}

	// Should have 2 triangles (duplicates allowed)
	if m.NumTriangles() != 2 {
		t.Fatalf("expected 2 triangles (duplicates), got %d", m.NumTriangles())
	}
}

func TestOverlapTriangleProhibited(t *testing.T) {
	// With WithOverlapTriangle(false), overlapping triangles should error
	m := NewMesh(WithOverlapTriangle(false))

	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{10, 0})
	v2, _ := m.AddVertex(types.Point{5, 10})

	_, err := m.AddPerimeter([]types.Point{{0, 0}, {10, 0}, {5, 10}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add triangle with one vertex order
	if err := m.AddTriangle(v0, v1, v2); err != nil {
		t.Fatalf("unexpected error adding first triangle: %v", err)
	}

	// Try to add same triangle with rotated vertex order (should fail)
	if err := m.AddTriangle(v1, v2, v0); err != ErrDuplicateTriangle {
		t.Fatalf("expected duplicate triangle error, got %v", err)
	}

	// Should have only 1 triangle
	if m.NumTriangles() != 1 {
		t.Fatalf("expected 1 triangle, got %d", m.NumTriangles())
	}
}

func TestOverlapTriangleAllowed(t *testing.T) {
	// With WithOverlapTriangle(true), overlapping triangles should be explicitly allowed
	m := NewMesh(WithOverlapTriangle(true))

	v0, _ := m.AddVertex(types.Point{0, 0})
	v1, _ := m.AddVertex(types.Point{10, 0})
	v2, _ := m.AddVertex(types.Point{5, 10})

	_, err := m.AddPerimeter([]types.Point{{0, 0}, {10, 0}, {5, 10}})
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Add triangle with one vertex order
	if err := m.AddTriangle(v0, v1, v2); err != nil {
		t.Fatalf("unexpected error adding first triangle: %v", err)
	}

	// Add same triangle with rotated vertex order (should succeed)
	if err := m.AddTriangle(v1, v2, v0); err != nil {
		t.Fatalf("expected overlapping triangle to succeed, got %v", err)
	}

	// Should have 2 triangles (duplicates allowed)
	if m.NumTriangles() != 2 {
		t.Fatalf("expected 2 triangles (duplicates), got %d", m.NumTriangles())
	}
}

func TestTriangleSpanningConcavePerimeter(t *testing.T) {
	// This test reproduces a bug where a triangle can span across a concave
	// perimeter section, with part inside and part outside the perimeter
	m := NewMesh(
		WithEpsilon(1e-9),
		WithMergeVertices(true),
		WithEdgeIntersectionCheck(true),
		WithTriangleEnforceNoVertexInside(true),
		WithEdgeCannotCrossPerimeter(true),
		WithOverlapTriangle(false),
	)

	perimeter := []types.Point{
		{X: 1.000, Y: 1.000},   // 0: bottom left
		{X: 54.000, Y: 1.000},  // 1: bottom right
		{X: 54.000, Y: 9.000},  // 2: top right
		{X: 50.000, Y: 9.000},  // 3: start of concave indent
		{X: 50.000, Y: 7.000},  // 4: down into indent
		{X: 49.000, Y: 6.000},  // 5: bottom of indent
		{X: 40.000, Y: 6.000},  // 6: bottom of indent
		{X: 39.000, Y: 7.000},  // 7: back up from indent
		{X: 39.000, Y: 9.000},  // 8: end of concave indent
		{X: 1.000, Y: 9.000},   // 9: top left
	}

	_, err := m.AddPerimeter(perimeter)
	if err != nil {
		t.Fatalf("failed to add perimeter: %v", err)
	}

	// Triangle 2,4,8 spans across the concave section at the top
	// Vertices:
	//   2: (54, 9) - top right corner
	//   4: (50, 7) - right side of concave section
	//   8: (39, 9) - left side of concave section
	//
	// The edge from vertex 8 (39,9) to vertex 2 (54,9) is a horizontal line
	// at y=9 that spans across the concave indent (which dips down to y=6).
	// The midpoint of this edge is at (46.5, 9), which is OUTSIDE the
	// perimeter because the perimeter dips down in that region.
	//
	// This triangle should be REJECTED because part of it extends outside
	// the perimeter boundary.
	err = m.AddTriangle(2, 4, 8)
	if err == nil {
		t.Error("Expected error when adding triangle that spans outside concave perimeter, got nil")
	} else if err != ErrEdgeCrossesPerimeter {
		t.Logf("Got error (expected ErrEdgeCrossesPerimeter): %v", err)
	}

	// Verify no triangles were added
	if m.NumTriangles() != 0 {
		t.Errorf("Expected 0 triangles after failed add, got %d", m.NumTriangles())
	}
}
