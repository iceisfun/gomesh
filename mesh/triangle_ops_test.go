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
	m := NewMesh(WithEdgeCannotCrossPerimeter(true))

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
