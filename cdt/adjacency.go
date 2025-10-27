package cdt

import (
	"fmt"

	"github.com/iceisfun/gomesh/algorithm/robust"
	"github.com/iceisfun/gomesh/types"
)

// TriSoup is a lightweight triangulation workspace with adjacency tracking.
type TriSoup struct {
	V   []types.Point // All vertices (including super-triangle vertices)
	Tri []Tri         // All triangles (some may be marked deleted)

	// Tracks which triangles use each edge (for O(1) neighbor lookup)
	edge2tri map[EdgeKey][]edgeUse

	// Free list for deleted triangle slots
	freeList []TriID
}

// edgeUse tracks a triangle that uses a particular edge.
type edgeUse struct {
	T         TriID
	LocalEdge int // Which edge of the triangle (0, 1, or 2)
}

// NewTriSoup creates a new empty triangulation workspace.
func NewTriSoup(pts []types.Point, reserveTris int) *TriSoup {
	return &TriSoup{
		V:        pts,
		Tri:      make([]Tri, 0, reserveTris),
		edge2tri: make(map[EdgeKey][]edgeUse),
		freeList: nil,
	}
}

// AddTri adds a triangle with vertices (a, b, c) and returns its ID.
// The triangle is added with no neighbors initially.
func (ts *TriSoup) AddTri(a, b, c int) TriID {
	tri := Tri{
		V: [3]int{a, b, c},
		N: [3]TriID{NilTri, NilTri, NilTri},
	}

	var id TriID
	if len(ts.freeList) > 0 {
		// Reuse a deleted slot
		id = ts.freeList[len(ts.freeList)-1]
		ts.freeList = ts.freeList[:len(ts.freeList)-1]
		ts.Tri[id] = tri
	} else {
		// Append new triangle
		id = TriID(len(ts.Tri))
		ts.Tri = append(ts.Tri, tri)
	}

	// Register edges
	ts.registerTriEdges(id)

	return id
}

// RemoveTri marks a triangle as deleted and unregisters its edges.
func (ts *TriSoup) RemoveTri(t TriID) {
	if t < 0 || int(t) >= len(ts.Tri) {
		return
	}

	// Update all neighbors to not point to this triangle anymore
	for i := 0; i < 3; i++ {
		neighbor := ts.Tri[t].N[i]
		if neighbor != NilTri && !ts.IsDeleted(neighbor) {
			// Find and clear the reference from neighbor to this triangle
			for j := 0; j < 3; j++ {
				if ts.Tri[neighbor].N[j] == t {
					ts.Tri[neighbor].N[j] = NilTri
					break
				}
			}
		}
	}

	ts.unregisterTriEdges(t)

	// Mark as deleted by setting all vertices to -1
	ts.Tri[t].V = [3]int{-1, -1, -1}
	ts.Tri[t].N = [3]TriID{NilTri, NilTri, NilTri}

	// Add to free list
	ts.freeList = append(ts.freeList, t)
}

// IsDeleted returns true if the triangle has been deleted.
func (ts *TriSoup) IsDeleted(t TriID) bool {
	if t < 0 || int(t) >= len(ts.Tri) {
		return true
	}
	return ts.Tri[t].V[0] < 0
}

// SetNeighbors sets all three neighbors of a triangle at once.
func (ts *TriSoup) SetNeighbors(t TriID, n0, n1, n2 TriID) {
	if ts.IsDeleted(t) {
		return
	}
	ts.Tri[t].N[0] = n0
	ts.Tri[t].N[1] = n1
	ts.Tri[t].N[2] = n2
}

// FindTriEdge finds which local edge (0, 1, or 2) connects vertices a and b.
// Returns (localEdge, true) if found, or (-1, false) if not found.
func (ts *TriSoup) FindTriEdge(t TriID, a, b int) (int, bool) {
	if ts.IsDeleted(t) {
		return -1, false
	}

	tri := &ts.Tri[t]
	for i := 0; i < 3; i++ {
		v1, v2 := tri.Edge(i)
		if (v1 == a && v2 == b) || (v1 == b && v2 == a) {
			return i, true
		}
	}
	return -1, false
}

// FlipEdge performs an edge flip on the shared edge between two triangles.
// Returns the two new triangle IDs and true if successful.
// Returns (NilTri, NilTri, false) if the flip cannot be performed.
func (ts *TriSoup) FlipEdge(tLeft TriID, eLeft int) (TriID, TriID, bool) {
	if ts.IsDeleted(tLeft) {
		return NilTri, NilTri, false
	}

	tRight := ts.Tri[tLeft].N[eLeft]
	if tRight == NilTri || ts.IsDeleted(tRight) {
		return NilTri, NilTri, false
	}

	// Find the corresponding edge in the right triangle
	leftV1, leftV2 := ts.Tri[tLeft].Edge(eLeft)
	eRight, ok := ts.FindTriEdge(tRight, leftV1, leftV2)
	if !ok {
		return NilTri, NilTri, false
	}

	// Get the four vertices of the quad
	// Left triangle: (leftV1, leftV2, apex)
	// Right triangle: (leftV1, leftV2, opposite)
	apex := ts.Tri[tLeft].V[eLeft]
	opposite := ts.Tri[tRight].V[eRight]

	// Check if the flip creates valid triangles (not inverted)
	p1 := ts.V[leftV1]
	p2 := ts.V[leftV2]
	pApex := ts.V[apex]
	pOpp := ts.V[opposite]

	// New triangles would be (apex, opposite, leftV2) and (opposite, apex, leftV1)
	if robust.Orient2D(pApex, pOpp, p2) <= 0 {
		return NilTri, NilTri, false
	}
	if robust.Orient2D(pOpp, pApex, p1) <= 0 {
		return NilTri, NilTri, false
	}

	leftVerts := ts.Tri[tLeft].V
	leftN := ts.Tri[tLeft].N
	rightVerts := ts.Tri[tRight].V
	rightN := ts.Tri[tRight].N

	ts.RemoveTri(tLeft)
	ts.RemoveTri(tRight)

	newLeft := addTriCCW(ts, apex, opposite, leftV2)
	newRight := addTriCCW(ts, opposite, apex, leftV1)

	linkTrianglesOnEdge(ts, newLeft, newRight, apex, opposite)

	for _, idx := range []int{(eLeft + 1) % 3, (eLeft + 2) % 3} {
		neighbor := leftN[idx]
		va := leftVerts[(idx+1)%3]
		vb := leftVerts[(idx+2)%3]
		if triangleHasEdge(ts, newLeft, va, vb) {
			attachNeighbor(ts, newLeft, va, vb, neighbor)
		} else {
			attachNeighbor(ts, newRight, va, vb, neighbor)
		}
	}

	for _, idx := range []int{(eRight + 1) % 3, (eRight + 2) % 3} {
		neighbor := rightN[idx]
		va := rightVerts[(idx+1)%3]
		vb := rightVerts[(idx+2)%3]
		if triangleHasEdge(ts, newLeft, va, vb) {
			attachNeighbor(ts, newLeft, va, vb, neighbor)
		} else {
			attachNeighbor(ts, newRight, va, vb, neighbor)
		}
	}

	return newLeft, newRight, true
}

// updateNeighborReference updates a neighbor's reference from oldT to newT.
func (ts *TriSoup) updateNeighborReference(neighbor, oldT, newT TriID) {
	if neighbor == NilTri || ts.IsDeleted(neighbor) {
		return
	}

	for i := 0; i < 3; i++ {
		if ts.Tri[neighbor].N[i] == oldT {
			ts.Tri[neighbor].N[i] = newT
			return
		}
	}
}

// registerTriEdges adds the triangle's edges to the edge index.
func (ts *TriSoup) registerTriEdges(t TriID) {
	tri := &ts.Tri[t]
	for i := 0; i < 3; i++ {
		v1, v2 := tri.Edge(i)
		key := NewEdgeKey(v1, v2)
		ts.edge2tri[key] = append(ts.edge2tri[key], edgeUse{T: t, LocalEdge: i})
	}
}

// unregisterTriEdges removes the triangle's edges from the edge index.
func (ts *TriSoup) unregisterTriEdges(t TriID) {
	tri := &ts.Tri[t]
	for i := 0; i < 3; i++ {
		v1, v2 := tri.Edge(i)
		key := NewEdgeKey(v1, v2)

		uses := ts.edge2tri[key]
		for j, use := range uses {
			if use.T == t {
				ts.edge2tri[key] = append(uses[:j], uses[j+1:]...)
				break
			}
		}

		if len(ts.edge2tri[key]) == 0 {
			delete(ts.edge2tri, key)
		}
	}
}

// FindEdgeTriangles returns the triangles that share the edge (a, b).
func (ts *TriSoup) FindEdgeTriangles(a, b int) []edgeUse {
	key := NewEdgeKey(a, b)
	return ts.edge2tri[key]
}

// Validate performs sanity checks on the triangulation.
func (ts *TriSoup) Validate() error {
	for i, tri := range ts.Tri {
		if tri.V[0] < 0 {
			continue // Deleted
		}

		// Check vertex indices are valid
		for j := 0; j < 3; j++ {
			if tri.V[j] < 0 || tri.V[j] >= len(ts.V) {
				return fmt.Errorf("triangle %d has invalid vertex index %d", i, tri.V[j])
			}
		}

		// Check neighbor symmetry
		for j := 0; j < 3; j++ {
			neighbor := tri.N[j]
			if neighbor == NilTri {
				continue
			}

			if ts.IsDeleted(neighbor) {
				return fmt.Errorf("triangle %d references deleted neighbor %d", i, neighbor)
			}

			// Find the shared edge
			v1, v2 := tri.Edge(j)
			eNeighbor, ok := ts.FindTriEdge(neighbor, v1, v2)
			if !ok {
				return fmt.Errorf("triangle %d and neighbor %d don't share expected edge", i, neighbor)
			}

			// Check that the neighbor points back
			if ts.Tri[neighbor].N[eNeighbor] != TriID(i) {
				return fmt.Errorf("neighbor symmetry broken between triangles %d and %d", i, neighbor)
			}
		}
	}

	return nil
}
