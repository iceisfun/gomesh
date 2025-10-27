package cdt

// TriID is a unique identifier for a triangle in the TriSoup.
type TriID int

const (
	// NilTri represents an invalid or missing triangle reference.
	NilTri TriID = -1
)

// EdgeKey represents an undirected edge between two vertex indices.
// The indices are always ordered (min, max) for consistent hashing.
type EdgeKey struct {
	A, B int
}

// NewEdgeKey constructs an EdgeKey with ordered vertex indices.
func NewEdgeKey(a, b int) EdgeKey {
	if a > b {
		a, b = b, a
	}
	return EdgeKey{A: a, B: b}
}

// Tri represents a triangle with three vertex indices and three neighbor references.
// N[i] is the neighbor across the edge opposite to vertex V[i].
type Tri struct {
	V [3]int   // Vertex indices into TriSoup.V
	N [3]TriID // Neighbor triangle IDs (NilTri if boundary)
}

// Edge returns the local edge index (0, 1, or 2) and its two vertex indices.
// Edge i is opposite to vertex V[i], connecting V[(i+1)%3] to V[(i+2)%3].
func (t *Tri) Edge(localEdge int) (int, int) {
	return t.V[(localEdge+1)%3], t.V[(localEdge+2)%3]
}

// Location describes where a point is located within the triangulation.
type Location struct {
	T      TriID // The containing triangle
	OnEdge bool  // True if the point lies on an edge
	Edge   int   // If OnEdge, which local edge (0, 1, or 2)
}
