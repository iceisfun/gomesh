package cdt

import (
	"fmt"

	"github.com/iceisfun/gomesh/algorithm/robust"
)

// InsertPoint inserts a vertex into the triangulation at the given location.
// Returns the IDs of the new triangles created and the edges that need legalization.
func InsertPoint(ts *TriSoup, loc Location, vidx int) ([]TriID, []EdgeToLegalize, error) {
	if ts.IsDeleted(loc.T) {
		return nil, nil, fmt.Errorf("cannot insert into deleted triangle")
	}

	if vidx < 0 || vidx >= len(ts.V) {
		return nil, nil, fmt.Errorf("invalid vertex index %d", vidx)
	}

	if loc.OnEdge {
		return insertPointOnEdge(ts, loc.T, loc.Edge, vidx)
	}
	return insertPointInTriangle(ts, loc.T, vidx)
}

// EdgeToLegalize represents an edge that may need to be flipped.
type EdgeToLegalize struct {
	T TriID // Triangle containing the edge
	E int   // Local edge index (0, 1, or 2)
}

// insertPointInTriangle splits triangle t by inserting vertex vidx.
// Creates three new triangles: (v0, v1, vidx), (v1, v2, vidx), (v2, v0, vidx).
func insertPointInTriangle(ts *TriSoup, t TriID, vidx int) ([]TriID, []EdgeToLegalize, error) {
	tri := ts.Tri[t]
	v0, v1, v2 := tri.V[0], tri.V[1], tri.V[2]
	n0, n1, n2 := tri.N[0], tri.N[1], tri.N[2]

	ts.RemoveTri(t)

	// Create three CCW triangles around the inserted vertex.
	t0 := addTriCCW(ts, v0, v1, vidx)
	t1 := addTriCCW(ts, v1, v2, vidx)
	t2 := addTriCCW(ts, v2, v0, vidx)

	linkTrianglesOnEdge(ts, t0, t1, v1, vidx)
	linkTrianglesOnEdge(ts, t0, t2, vidx, v0)
	linkTrianglesOnEdge(ts, t1, t2, v2, vidx)

	attachNeighbor(ts, t1, v1, v2, n0)
	attachNeighbor(ts, t2, v2, v0, n1)
	attachNeighbor(ts, t0, v0, v1, n2)

	return []TriID{t0, t1, t2}, []EdgeToLegalize{
		{T: t0, E: 2},
		{T: t1, E: 2},
		{T: t2, E: 2},
	}, nil
}

// insertPointOnEdge splits the two triangles sharing edge e of triangle t.
// Creates four new triangles.
func insertPointOnEdge(ts *TriSoup, t TriID, e int, vidx int) ([]TriID, []EdgeToLegalize, error) {
	tri := ts.Tri[t]
	tOpp := tri.N[e]
	if tOpp == NilTri {
		return insertPointOnBoundaryEdge(ts, t, e, vidx)
	}
	if ts.IsDeleted(tOpp) {
		return nil, nil, fmt.Errorf("neighbor triangle is deleted")
	}

	v1, v2 := tri.Edge(e)
	eOpp, ok := ts.FindTriEdge(tOpp, v1, v2)
	if !ok {
		return nil, nil, fmt.Errorf("could not find shared edge in neighbor")
	}

	v0 := tri.V[e]
	vOpp := ts.Tri[tOpp].V[eOpp]

	nLeft1 := tri.N[(e+1)%3]
	nLeft2 := tri.N[(e+2)%3]
	nRight1 := ts.Tri[tOpp].N[(eOpp+1)%3]
	nRight2 := ts.Tri[tOpp].N[(eOpp+2)%3]

	ts.RemoveTri(t)
	ts.RemoveTri(tOpp)

	// Maintain winding so vidx is last when possible.
	t0 := addTriCCW(ts, v0, v1, vidx)
	t1 := addTriCCW(ts, v0, vidx, v2)
	t2 := addTriCCW(ts, vOpp, vidx, v1)
	t3 := addTriCCW(ts, vOpp, v2, vidx)

	ts.Tri[t0].N[0] = t2 // (v1, vidx)
	ts.Tri[t0].N[1] = t1 // (vidx, v0)
	ts.Tri[t0].N[2] = nLeft2

	ts.Tri[t1].N[0] = t3 // (vidx, v2)
	ts.Tri[t1].N[1] = nLeft1
	ts.Tri[t1].N[2] = t0 // (v0, vidx)

	ts.Tri[t2].N[0] = t0 // (vidx, v1)
	ts.Tri[t2].N[1] = nRight1
	ts.Tri[t2].N[2] = t3 // (vOpp, vidx)

	ts.Tri[t3].N[0] = t1 // (v2, vidx)
	ts.Tri[t3].N[1] = t2
	ts.Tri[t3].N[2] = nRight2

	linkTrianglesOnEdge(ts, t0, t1, v0, vidx)
	linkTrianglesOnEdge(ts, t0, t2, vidx, v1)
	linkTrianglesOnEdge(ts, t1, t3, vidx, v2)
	linkTrianglesOnEdge(ts, t2, t3, vOpp, vidx)

	attachNeighbor(ts, t0, v0, v1, nLeft2)
	attachNeighbor(ts, t1, v2, v0, nLeft1)
	attachNeighbor(ts, t2, v1, vOpp, nRight1)
	attachNeighbor(ts, t3, vOpp, v2, nRight2)

	return []TriID{t0, t1, t2, t3}, []EdgeToLegalize{
		{T: t0, E: 2}, // (v0, v1)
		{T: t1, E: 1}, // (v2, v0)
		{T: t2, E: 1}, // (v1, vOpp)
		{T: t3, E: 2}, // (vOpp, v2)
	}, nil
}

// insertPointOnBoundaryEdge handles inserting a point on an edge that has no neighbor.
func insertPointOnBoundaryEdge(ts *TriSoup, t TriID, e int, vidx int) ([]TriID, []EdgeToLegalize, error) {
	tri := ts.Tri[t]
	v0 := tri.V[e]
	v1, v2 := tri.Edge(e)

	// Save neighbors
	n1 := tri.N[(e+1)%3]
	n2 := tri.N[(e+2)%3]

	// Remove old triangle
	ts.RemoveTri(t)

	// Create two new triangles with winding matching the exterior fan
	t0 := addTriCCW(ts, v0, v1, vidx)
	t1 := addTriCCW(ts, v0, vidx, v2)

	ts.Tri[t0].N[0] = NilTri
	ts.Tri[t0].N[1] = t1
	ts.Tri[t0].N[2] = n2

	ts.Tri[t1].N[0] = NilTri
	ts.Tri[t1].N[1] = n1
	ts.Tri[t1].N[2] = t0

	// Update external neighbors
	linkTrianglesOnEdge(ts, t0, t1, v0, vidx)
	attachNeighbor(ts, t1, v2, v0, n1)
	attachNeighbor(ts, t0, v0, v1, n2)

	edgesToLegalize := []EdgeToLegalize{
		{T: t0, E: 2}, // Edge (v0, v1) - external
		{T: t1, E: 1}, // Edge (v2, v0) - external
	}

	return []TriID{t0, t1}, edgesToLegalize, nil
}

func addTriCCW(ts *TriSoup, a, b, c int) TriID {
	pa, pb, pc := ts.V[a], ts.V[b], ts.V[c]
	if robust.Orient2D(pa, pb, pc) < 0 {
		b, c = c, b
	}
	return ts.AddTri(a, b, c)
}

func linkTrianglesOnEdge(ts *TriSoup, tA, tB TriID, a, b int) {
	if ts.IsDeleted(tA) || ts.IsDeleted(tB) {
		return
	}
	edgeA, okA := ts.FindTriEdge(tA, a, b)
	edgeB, okB := ts.FindTriEdge(tB, a, b)
	if !okA || !okB {
		return
	}
	ts.Tri[tA].N[edgeA] = tB
	ts.Tri[tB].N[edgeB] = tA
}

func attachNeighbor(ts *TriSoup, t TriID, a, b int, neighbor TriID) {
	edgeIdx, ok := ts.FindTriEdge(t, a, b)
	if !ok {
		return
	}
	ts.Tri[t].N[edgeIdx] = neighbor
	if neighbor == NilTri || ts.IsDeleted(neighbor) {
		return
	}
	if nEdge, ok := ts.FindTriEdge(neighbor, a, b); ok {
		ts.Tri[neighbor].N[nEdge] = t
	}
}

func triangleHasEdge(ts *TriSoup, t TriID, a, b int) bool {
	_, ok := ts.FindTriEdge(t, a, b)
	return ok
}
