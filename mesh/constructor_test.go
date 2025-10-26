package mesh

import "testing"

func TestNewMeshDefaults(t *testing.T) {
	m := NewMesh()
	if m == nil {
		t.Fatalf("expected mesh instance")
	}
	if m.NumVertices() != 0 || m.NumTriangles() != 0 {
		t.Fatalf("expected empty mesh")
	}
	if m.vertexIndex != nil {
		t.Fatalf("vertex index should be nil without merging")
	}
}

func TestNewMeshWithOptions(t *testing.T) {
	m := NewMesh(WithMergeVertices(true), WithMergeDistance(0.1))
	if m.vertexIndex == nil {
		t.Fatalf("expected vertex index when merging enabled")
	}
}
