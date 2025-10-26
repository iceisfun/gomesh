package types

// Edge represents an undirected connection between two vertices.
//
// Edges are stored in canonical form with vertex IDs in ascending order,
// ensuring that Edge{a, b} and Edge{b, a} compare as equal.
//
// Use NewEdge() to construct edges in canonical form, or use Canonical()
// to normalize an existing edge.
//
// Example:
//
//	e1 := types.NewEdge(5, 3)  // Stored as Edge{3, 5}
//	e2 := types.NewEdge(3, 5)  // Stored as Edge{3, 5}
//	// e1 == e2 (true)
type Edge [2]VertexID

// NewEdge creates an edge in canonical form (min ID first).
func NewEdge(v1, v2 VertexID) Edge {
	if v1 < v2 {
		return Edge{v1, v2}
	}
	return Edge{v2, v1}
}

// Canonical returns this edge in canonical form.
func (e Edge) Canonical() Edge {
	return NewEdge(e[0], e[1])
}

// IsCanonical returns true if this edge is in canonical form.
func (e Edge) IsCanonical() bool {
	return e[0] <= e[1]
}

// V1 returns the first vertex ID (smaller ID in canonical form).
func (e Edge) V1() VertexID {
	return e[0]
}

// V2 returns the second vertex ID (larger ID in canonical form).
func (e Edge) V2() VertexID {
	return e[1]
}
