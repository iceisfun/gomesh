package types

import "testing"

func TestTriangleAccessors(t *testing.T) {
	tri := NewTriangle(1, 2, 3)
	if tri.V1() != 1 || tri.V2() != 2 || tri.V3() != 3 {
		t.Fatalf("unexpected triangle vertices: %+v", tri)
	}
	edges := tri.Edges()
	expected := [3]Edge{NewEdge(1, 2), NewEdge(2, 3), NewEdge(1, 3)}
	if edges != expected {
		t.Fatalf("unexpected edges: got %v expected %v", edges, expected)
	}
}

func TestTriangleVerticesSlice(t *testing.T) {
	tri := Triangle{5, 6, 7}
	verts := tri.Vertices()
	if len(verts) != 3 || verts[0] != 5 || verts[1] != 6 || verts[2] != 7 {
		t.Fatalf("unexpected vertices slice: %v", verts)
	}
}
