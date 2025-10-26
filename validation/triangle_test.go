package validation

import (
	"errors"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

type mockMesh struct {
	vertices  []types.Point
	edgeSet   map[types.Edge]struct{}
	triangles map[[3]types.VertexID]types.Triangle
}

func newMockMesh(points []types.Point) *mockMesh {
	return &mockMesh{
		vertices:  points,
		edgeSet:   make(map[types.Edge]struct{}),
		triangles: make(map[[3]types.VertexID]types.Triangle),
	}
}

func (m *mockMesh) NumVertices() int { return len(m.vertices) }

func (m *mockMesh) GetVertex(id types.VertexID) types.Point { return m.vertices[id] }

func (m *mockMesh) EdgeSet() map[types.Edge]struct{} { return m.edgeSet }

func (m *mockMesh) HasTriangleWithKey(key [3]types.VertexID) (types.Triangle, bool) {
	t, ok := m.triangles[key]
	return t, ok
}

func (m *mockMesh) addTriangle(tri types.Triangle) {
	key := CanonicalTriangleKey(tri)
	m.triangles[key] = tri

	for _, edge := range tri.Edges() {
		m.edgeSet[edge] = struct{}{}
	}
}

func TestValidateTriangleDegenerate(t *testing.T) {
	mesh := newMockMesh([]types.Point{{0, 0}, {1, 1}, {2, 2}})
	cfg := Config{Epsilon: 1e-9}
	err := ValidateTriangle(types.Triangle{0, 1, 2}, mesh.vertices[0], mesh.vertices[1], mesh.vertices[2], cfg, mesh)
	if !errors.Is(err, Errors().Degenerate) {
		t.Fatalf("expected degenerate error, got %v", err)
	}
}

func TestValidateTriangleDuplicate(t *testing.T) {
	mesh := newMockMesh([]types.Point{{0, 0}, {1, 0}, {0, 1}})
	tri := types.Triangle{0, 1, 2}
	mesh.addTriangle(tri)

	cfg := Config{Epsilon: 1e-9, ErrorOnDuplicateTriangle: true}
	err := ValidateTriangle(tri, mesh.vertices[0], mesh.vertices[1], mesh.vertices[2], cfg, mesh)
	if !errors.Is(err, Errors().Duplicate) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestValidateTriangleOpposingDuplicate(t *testing.T) {
	mesh := newMockMesh([]types.Point{{0, 0}, {1, 0}, {0, 1}})
	tri := types.Triangle{0, 1, 2}
	mesh.addTriangle(tri)

	cfg := Config{Epsilon: 1e-9, ErrorOnOpposingDuplicate: true}
	err := ValidateTriangle(types.Triangle{0, 2, 1}, mesh.vertices[0], mesh.vertices[2], mesh.vertices[1], cfg, mesh)
	if !errors.Is(err, Errors().OpposingDuplicate) {
		t.Fatalf("expected opposing duplicate error, got %v", err)
	}
}

func TestValidateTriangleVertexInside(t *testing.T) {
	mesh := newMockMesh([]types.Point{{0, 0}, {2, 0}, {0, 2}, {0.5, 0.5}})
	tri := types.Triangle{0, 1, 2}
	cfg := Config{Epsilon: 1e-9, ValidateVertexInside: true}
	err := ValidateTriangle(tri, mesh.vertices[0], mesh.vertices[1], mesh.vertices[2], cfg, mesh)
	if !errors.Is(err, Errors().VertexInside) {
		t.Fatalf("expected vertex-inside error, got %v", err)
	}
}
