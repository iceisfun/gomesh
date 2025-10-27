package cdt

import (
	"github.com/iceisfun/gomesh/algorithm/robust"
)

// IsEdgeConstrained checks if an edge is marked as constrained.
func IsEdgeConstrained(edge EdgeKey, constrained map[EdgeKey]bool) bool {
	return constrained[edge]
}

// IsIllegal checks if an edge between two triangles violates the Delaunay property.
// An edge is illegal if:
//  1. It is not a constrained edge, AND
//  2. The opposite vertex in the neighboring triangle is strictly inside the circumcircle
//     of the triangle containing the edge.
//
// Returns true if the edge should be flipped.
func IsIllegal(ts *TriSoup, t TriID, e int, constrained map[EdgeKey]bool) bool {
	if ts.IsDeleted(t) {
		return false
	}

	tri := &ts.Tri[t]
	neighbor := tri.N[e]

	// Boundary edges can't be illegal
	if neighbor == NilTri || ts.IsDeleted(neighbor) {
		return false
	}

	// Get edge vertices
	v1, v2 := tri.Edge(e)
	edgeKey := NewEdgeKey(v1, v2)

	// Constrained edges are never illegal
	if IsEdgeConstrained(edgeKey, constrained) {
		return false
	}

	// Find the opposite vertex in the neighboring triangle
	triOpp := &ts.Tri[neighbor]
	eOpp, ok := ts.FindTriEdge(neighbor, v1, v2)
	if !ok {
		return false
	}
	vOppIdx := triOpp.V[eOpp]
	d := ts.V[vOppIdx]

	// Use InCircle predicate to check if d is inside the circumcircle of (a, b, c)
	// The triangle vertices should be in CCW order for the predicate
	// We need to ensure we're testing the correct orientation

	// Get the vertex opposite to edge e in triangle t
	apex := tri.V[e]
	apexPt := ts.V[apex]

	// The edge connects v1 and v2
	p1 := ts.V[v1]
	p2 := ts.V[v2]

	// Check if d is in the circumcircle of the triangle (apexPt, p1, p2)
	// Make sure the triangle is CCW
	orient := robust.Orient2D(apexPt, p1, p2)
	var inCircle int
	if orient > 0 {
		// CCW
		inCircle = robust.InCircle(apexPt, p1, p2, d)
	} else if orient < 0 {
		// CW - reverse order
		inCircle = robust.InCircle(apexPt, p2, p1, d)
	} else {
		// Degenerate triangle
		return false
	}

	// If d is strictly inside the circumcircle, the edge is illegal
	return inCircle > 0
}

// LegalizeAround performs edge legalization starting from a set of seed edges.
// It uses a queue to process edges that might be illegal and flips them if needed.
// This continues until all edges satisfy the Delaunay property (or are constrained).
func LegalizeAround(ts *TriSoup, seeds []EdgeToLegalize, constrained map[EdgeKey]bool) {
	if constrained == nil {
		constrained = make(map[EdgeKey]bool)
	}

	// Use a queue for BFS-style legalization
	queue := make([]EdgeToLegalize, len(seeds))
	copy(queue, seeds)

	// Track edges we've already processed to avoid infinite loops
	processed := make(map[edgeRef]bool)

	for len(queue) > 0 {
		// Pop from queue
		edge := queue[0]
		queue = queue[1:]

		// Skip if already processed
		ref := edgeRef{T: edge.T, E: edge.E}
		if processed[ref] {
			continue
		}
		processed[ref] = true

		// Skip if triangle was deleted
		if ts.IsDeleted(edge.T) {
			continue
		}

		// Check if the edge is illegal
		if !IsIllegal(ts, edge.T, edge.E, constrained) {
			continue
		}

		// Perform the flip
		newLeft, newRight, ok := ts.FlipEdge(edge.T, edge.E)
		if !ok {
			continue
		}

		// Add the four edges around the new diamond to the queue
		// These are the edges that might have become illegal due to the flip

		// For newLeft triangle, check edges 0 and 2 (not the shared edge 1)
		queue = append(queue, EdgeToLegalize{T: newLeft, E: 0})
		queue = append(queue, EdgeToLegalize{T: newLeft, E: 2})

		// For newRight triangle, check edges 0 and 2 (not the shared edge 1)
		queue = append(queue, EdgeToLegalize{T: newRight, E: 0})
		queue = append(queue, EdgeToLegalize{T: newRight, E: 2})
	}
}

// edgeRef uniquely identifies an edge within the triangulation.
type edgeRef struct {
	T TriID
	E int
}

// LegalizeEdge is a simpler interface that legalizes a single edge recursively.
// This is useful for testing or when you want to legalize a specific edge.
func LegalizeEdge(ts *TriSoup, t TriID, e int, constrained map[EdgeKey]bool) {
	LegalizeAround(ts, []EdgeToLegalize{{T: t, E: e}}, constrained)
}

// IsDelaunay checks if the entire triangulation satisfies the Delaunay property
// (ignoring constrained edges). This is useful for validation and testing.
func IsDelaunay(ts *TriSoup, constrained map[EdgeKey]bool) bool {
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		for e := 0; e < 3; e++ {
			if IsIllegal(ts, TriID(i), e, constrained) {
				return false
			}
		}
	}
	return true
}
