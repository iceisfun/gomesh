package mesh

import (
	"github.com/iceisfun/gomesh/spatial"
	"github.com/iceisfun/gomesh/types"
)

// Mesh represents a 2D triangle mesh with validated topology.
type Mesh struct {
	vertices  []types.Point
	triangles []types.Triangle

	cfg config

	vertexIndex spatial.Index

	edgeSet map[types.Edge]struct{}

	triangleSet map[[3]types.VertexID]types.Triangle

	perimeters []types.PolygonLoop
	holes      []types.PolygonLoop
}

// NumVertices returns the number of vertices in the mesh.
func (m *Mesh) NumVertices() int {
	return len(m.vertices)
}

// NumTriangles returns the number of triangles in the mesh.
func (m *Mesh) NumTriangles() int {
	return len(m.triangles)
}

// GetVertex returns the coordinates of a vertex by ID.
func (m *Mesh) GetVertex(id types.VertexID) types.Point {
	return m.vertices[id]
}

// GetTriangle returns a triangle by index.
func (m *Mesh) GetTriangle(idx int) types.Triangle {
	return m.triangles[idx]
}

// GetVertices returns a copy of all vertex coordinates.
func (m *Mesh) GetVertices() []types.Point {
	out := make([]types.Point, len(m.vertices))
	copy(out, m.vertices)
	return out
}

// GetTriangles returns a copy of all triangles.
func (m *Mesh) GetTriangles() []types.Triangle {
	out := make([]types.Triangle, len(m.triangles))
	copy(out, m.triangles)
	return out
}

// GetTriangleCoords returns the coordinates of a triangle's vertices.
func (m *Mesh) GetTriangleCoords(idx int) (types.Point, types.Point, types.Point) {
	t := m.triangles[idx]
	return m.vertices[t.V1()], m.vertices[t.V2()], m.vertices[t.V3()]
}

// IsValidVertexID reports whether the supplied ID references an existing vertex.
func (m *Mesh) IsValidVertexID(id types.VertexID) bool {
	return id >= 0 && int(id) < len(m.vertices)
}

// Epsilon returns the configured epsilon tolerance.
func (m *Mesh) Epsilon() float64 {
	return m.cfg.epsilon
}

// EdgeSet exposes the set of edges currently tracked by the mesh.
func (m *Mesh) EdgeSet() map[types.Edge]struct{} {
	return m.edgeSet
}

// HasTriangleWithKey reports whether the canonical key is present.
func (m *Mesh) HasTriangleWithKey(key [3]types.VertexID) (types.Triangle, bool) {
	tri, ok := m.triangleSet[key]
	return tri, ok
}
