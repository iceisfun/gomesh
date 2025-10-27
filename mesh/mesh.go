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

// EdgeUsageCounts returns a map of each edge to the number of triangles using it.
//
// In a valid triangulation, each edge should be used by at most 2 triangles.
func (m *Mesh) EdgeUsageCounts() map[types.Edge]int {
	counts := make(map[types.Edge]int)
	for _, tri := range m.triangles {
		edges := tri.Edges()
		for _, edge := range edges {
			counts[edge]++
		}
	}
	return counts
}

// GetUntriangulatedVertices returns vertices from the given loops that are not part of any triangle.
//
// This is useful for identifying areas with missing triangulation during debugging.
//
// Example:
//
//	loops := []types.PolygonLoop{perimeter, hole1, hole2}
//	untriangulated := m.GetUntriangulatedVertices(loops)
//	if len(untriangulated) > 0 {
//	    fmt.Printf("Found %d untriangulated vertices\n", len(untriangulated))
//	}
func (m *Mesh) GetUntriangulatedVertices(loops []types.PolygonLoop) []types.VertexID {
	// Build set of all vertices used in triangles
	triangulatedVertices := make(map[types.VertexID]struct{})
	for _, tri := range m.triangles {
		triangulatedVertices[tri.V1()] = struct{}{}
		triangulatedVertices[tri.V2()] = struct{}{}
		triangulatedVertices[tri.V3()] = struct{}{}
	}

	// Collect unique vertices from loops
	loopVertices := make(map[types.VertexID]struct{})
	for _, loop := range loops {
		for _, vid := range loop {
			loopVertices[vid] = struct{}{}
		}
	}

	// Find vertices in loops that aren't triangulated
	var untriangulated []types.VertexID
	for vid := range loopVertices {
		if _, ok := triangulatedVertices[vid]; !ok {
			untriangulated = append(untriangulated, vid)
		}
	}

	return untriangulated
}
