package mesh

import (
	"testing"

	"gomesh/types"
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
