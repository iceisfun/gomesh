package spatial

import "gomesh/types"

// Index provides spatial queries for vertices.
type Index interface {
	// FindVerticesNear returns vertex IDs within radius of point p.
	FindVerticesNear(p types.Point, radius float64) []types.VertexID
	// AddVertex adds a vertex to the index.
	AddVertex(id types.VertexID, p types.Point)
	// Build finalizes the index structure.
	Build()
}
