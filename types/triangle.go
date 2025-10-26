package types

// Triangle represents an ordered triplet of vertices forming a triangle.
//
// The order of vertices determines the winding direction:
//   - Counter-clockwise (CCW) order yields positive signed area
//   - Clockwise (CW) order yields negative signed area
//   - Collinear vertices yield zero (or near-zero) signed area
//
// Triangles are stored exactly as provided; no automatic reordering
// is performed. Use predicates.Area2 or predicates.Orient to determine
// winding.
//
// Example:
//
//	t := types.Triangle{0, 1, 2}  // CCW if vertices are positioned appropriately
type Triangle [3]VertexID

// NewTriangle creates a triangle from three vertex IDs.
func NewTriangle(v1, v2, v3 VertexID) Triangle {
	return Triangle{v1, v2, v3}
}

// V1 returns the first vertex.
func (t Triangle) V1() VertexID {
	return t[0]
}

// V2 returns the second vertex.
func (t Triangle) V2() VertexID {
	return t[1]
}

// V3 returns the third vertex.
func (t Triangle) V3() VertexID {
	return t[2]
}

// Vertices returns all three vertex IDs as a slice.
func (t Triangle) Vertices() []VertexID {
	return []VertexID{t[0], t[1], t[2]}
}

// Edges returns the three edges of this triangle in canonical form.
//
// The edges are returned in the order: (v1,v2), (v2,v3), (v3,v1).
func (t Triangle) Edges() [3]Edge {
	return [3]Edge{
		NewEdge(t[0], t[1]),
		NewEdge(t[1], t[2]),
		NewEdge(t[2], t[0]),
	}
}
